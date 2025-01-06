package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/captionschain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/clipstreamer"
	"gitlab.com/andyinabox/yesterdaysnews/domain/downloader"
	"gitlab.com/andyinabox/yesterdaysnews/domain/uploader"
	"gitlab.com/andyinabox/yesterdaysnews/domain/videoprocessor"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

var verbose bool
var timestamp string
var yesterday time.Time
var prefixLength int

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	log.SetReportCaller(true)
	log.SetReportTimestamp(false)

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	err := godotenv.Load()
	if err != nil {
		log.Warnf("error loading .env: %s", err)
	}

	yesterday = util.Yesterday()
	timestamp = util.Timestamp(yesterday)
	prefixLength = 2
}

func main() {
	var err error
	var eh errorhandler.ErrorHandler
	var dl domain.Downloader
	var vp domain.VideoProcessor
	var up domain.Uploader
	var cs domain.ClipStreamer
	var cc domain.CaptionsChain
	var mu sync.Mutex

	ctx := context.Background()

	manifest := &domain.Manifest{
		Date: yesterday,
		ID:   timestamp,
		Files: domain.ManifestFiles{
			Clips: []string{},
		},
	}

	fatalFunc := func(typ string, err error) {
		log.Fatalf("%s: %s", typ, err)
	}

	errorFund := func(typ string, err error) {
		if typ == domain.ErrTypeFatal {
			fatalFunc(typ, err)
		}
		log.Errorf("%s: %s", typ, err)
	}

	eh = errorhandler.New(ctx, &errorhandler.Config{
		ErrorFunc: errorFund,
		FatalFunc: fatalFunc,
		Thresholds: map[string]int{
			domain.ErrTypeBuildModel:     1,
			domain.ErrTypeSaveModel:      1,
			domain.ErrTypeUploadModel:    1,
			domain.ErrTypeSaveManifest:   1,
			domain.ErrTypeUploadManifest: 1,
		},
	})

	errs := eh.Channel()

	// error recovery
	defer func() {

		// print error report if there are any errors
		if eh.CountAll() > 0 {
			report := eh.Report()
			log.Print(eh.Report())
			_ = os.WriteFile(fmt.Sprintf("errors.%s.json", timestamp), []byte(report), os.ModePerm)
		}
	}()

	dl = downloader.New(&downloader.Config{
		GoogleAPIKey: os.Getenv("GOOGLE_API_KEY"),
	})

	vp = videoprocessor.New(&videoprocessor.Config{})

	up = uploader.New(&uploader.Config{
		S3Endpoint:  os.Getenv("YN_S3_ENDPOINT"),
		S3AccessKey: os.Getenv("YN_S3_ACCESS_KEY"),
		S3SecretKey: os.Getenv("YN_S3_SECRET_ACCESS_KEY"),

		ContainerName: os.Getenv("YN_S3_BUCKET_NAME"),
		PrimaryDir:    "current",
	})

	downloadDir := "download"
	err = os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	outputDir := "dist"
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	clipsDir := "clips"
	err = os.MkdirAll(filepath.Join(outputDir, clipsDir), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	cs = clipstreamer.New(dl, vp, up, eh.Channel(), &clipstreamer.Config{
		VideoDate:                yesterday,
		DownloadCountPerPlaylist: 1, //10,
		DownloadDir:              downloadDir,
		OutputDir:                outputDir,
		ClipsDir:                 clipsDir,
		MinClipLength:            5 * time.Second,
		MaxClipLength:            20 * time.Second,
		FileUploadDir:            timestamp,
	})

	// TODO: let's just use the playlist IDs and skip this step
	channelNames := []string{"@cnn", "@msnbc", "@foxnews"}
	playlistIDs := make([]string, len(channelNames))

	for i, name := range channelNames {
		playlistId, err := dl.GetChannelPlaylistID(ctx, name)
		if err != nil {
			log.Fatal(err)
		}
		playlistIDs[i] = playlistId
	}

	ids := cs.VideoIDStream(ctx, playlistIDs...)
	paths := cs.VideoDownloadStream(ctx, ids)
	clips := cs.VideoCutStream(ctx, paths)
	clipUploades := cs.VideoUploadStream(ctx, clips)

	for upload := range clipUploades {
		mu.Lock()
		manifest.Files.Clips = append(manifest.Files.Clips, strings.Replace(upload, timestamp+"/", "", 1))
		mu.Unlock()
	}

	cc = captionschain.New(prefixLength)
	cc.BuildFromMultiple([]domain.Corpus{
		{
			Type:     domain.CorpusTypeVTT,
			FileGlob: "download/*.vtt",
			Weight:   1,
		},
		{
			Type:     domain.CorpusTypeText,
			FileGlob: "hospital.txt",
			Weight:   3,
		},
	})

	data, err := cc.Save()
	if err != nil {
		errs <- errorhandler.Err(domain.ErrTypeBuildModel, err)
		return
	}

	modelFilePath := filepath.Join(outputDir, domain.ManifestModelFileName)
	log.Infof("saving model file to %q", modelFilePath)
	err = os.WriteFile(modelFilePath, data, os.ModePerm)
	if err != nil {
		errs <- errorhandler.Err(domain.ErrTypeSaveModel, err)
		return
	}

	modelFileKey := filepath.Join(timestamp, domain.ManifestModelFileName)
	log.Infof("uploading %q as %q", modelFilePath, modelFileKey)
	modelFileKey, err = up.UploadFile(ctx, modelFilePath, modelFileKey, "application/json", false)
	if err != nil {
		errs <- errorhandler.Err(domain.ErrTypeUploadModel, err)
		return
	}

	mu.Lock()
	manifest.Files.ModelFile = strings.Replace(modelFileKey, timestamp+"/", "", 1)
	mu.Unlock()

	log.Info("saving manifest")
	data, err = json.Marshal(manifest)
	if err != nil {
		errs <- errorhandler.Err(domain.ErrTypeSaveManifest, err)
		return
	}

	manifestFilePath := filepath.Join(outputDir, "manifest.json")
	err = os.WriteFile(manifestFilePath, data, os.ModePerm)
	if err != nil {
		errs <- errorhandler.Err(domain.ErrTypeSaveManifest, err)
		return
	}

	manifestFileKey := filepath.Join(timestamp, "manifest.json")
	log.Infof("uploading %q as %q", manifestFilePath, manifestFileKey)
	modelFileKey, err = up.UploadFile(ctx, manifestFilePath, manifestFileKey, "application/json", false)
	if err != nil {
		errs <- errorhandler.Err(domain.ErrTypeUploadManifest, err)
		return
	}

	log.Info("Done")

}
