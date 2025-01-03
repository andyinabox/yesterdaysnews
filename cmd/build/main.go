package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/clipstreamer"
	"gitlab.com/andyinabox/yesterdaysnews/domain/downloader"
	"gitlab.com/andyinabox/yesterdaysnews/domain/uploader"
	"gitlab.com/andyinabox/yesterdaysnews/domain/videoprocessor"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

var verbose bool

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

}

func main() {
	var err error
	var dl domain.Downloader
	var vp domain.VideoProcessor
	var up domain.Uploader
	var cs domain.ClipStreamer

	dl = downloader.New(&downloader.Config{
		GoogleAPIKey: os.Getenv("YN_GOOGLE_API_KEY"),
	})

	vp = videoprocessor.New(&videoprocessor.Config{})

	up = uploader.New(&uploader.Config{})

	downloadDir := "download"
	err = os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	outputDir := "output"
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	cs = clipstreamer.New(dl, vp, up, &clipstreamer.Config{
		VideoDate:                util.Yesterday(),
		DownloadCountPerPlaylist: 10,
		DownloadDir:              downloadDir,
		OutputDir:                outputDir,
		MinClipLength:            5 * time.Second,
		MaxClipLength:            20 * time.Second,
		FileUploadDir:            "test",
	})

	ctx, _ := context.WithCancelCause(context.Background())

	errHandler := func(err domain.StreamErr) {
		switch err.Type() {
		case domain.StreamErrFatal:
			log.Fatal(err)
		default:
			log.Errorf("%s: %s", err.Type(), err)
		}
	}

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

	errs := cs.ErrorStream(ctx, errHandler)
	ids := cs.VideoIDStream(ctx, errs, playlistIDs...)
	paths := cs.VideoDownloadStream(ctx, errs, ids)

	for path := range paths {
		log.Info(path)
	}

	log.Info("Done")

}
