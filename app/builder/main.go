package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/builder"
	"gitlab.com/andyinabox/yesterdaysnews/domain/captionschain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/domain/logger"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/configloader"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/streams"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

//go:embed hospital.txt
var hospitalText string

//go:embed overlay.png
var overlayImage []byte

const defaultPlaylists = "UUupvZG-5ko_eiXAupbDfxWw,UUaXkIU1QidjPwiAYu6GcHjg,UUXIJgqnII2ZOINSWNOGFThA"

var (
	verbose    bool
	loggerType string
	buildPhase string
	buildId    string

	playlistIDs              string
	objectStoreContainerName string
	outputDir                string
	throttleDownloadsBy      string

	maxVideoSize                                         int
	downloadCountPerPlaylist                             int
	maxPlaylistRequests                                  int
	minClipLengthSeconds, maxClipLengthSeconds           int
	captionPrefixLength                                  int
	captionNewsCorpusWeight, captionHospitalCorpusWeight int
	captionMinDuration, captionMaxDuration               float64
	totalBuildsToKeep                                    int

	keepOutputFiles, skipUpload bool
)

func init() {
	// meta flags
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&loggerType, "log", "text", "logger type (text, json, loki)")
	flag.StringVar(&buildPhase, "b", "all", "build phase to execute")
	flag.StringVar(&buildId, "id", "", "build ID")

	// config flags
	flag.StringVar(&playlistIDs, "playlistids", defaultPlaylists, "comma-separated list of playlists to download")
	flag.StringVar(&objectStoreContainerName, "containername", "yesterdaysnews", "object storage container name")
	flag.StringVar(&outputDir, "output", "dist", "dir to output build artifacts to")
	flag.StringVar(&throttleDownloadsBy, "throttledl", "0s", "throttle downloads by this amount")

	flag.IntVar(&maxVideoSize, "maxvideosize", 52428800, "max video download size in bytes")
	flag.IntVar(&downloadCountPerPlaylist, "count", 10, "download count per playlist")
	flag.IntVar(&maxPlaylistRequests, "maxplaylistreq", 3, "the maximum times to request a new playlist page before giving up")
	flag.IntVar(&minClipLengthSeconds, "mincliplength", 5, "minimum clip length in seconds")
	flag.IntVar(&maxClipLengthSeconds, "maxcliplength", 15, "maximum clip length in seconds")
	flag.IntVar(&captionPrefixLength, "prefixlength", 2, "caption chain prefix length")
	flag.IntVar(&captionNewsCorpusWeight, "newsweight", 1, "weight for the news corpus in chain")
	flag.IntVar(&captionHospitalCorpusWeight, "hospitalweight", 1, "weight for the hospital corpus in chain")
	flag.Float64Var(&captionMinDuration, "mincap", 3.0, "minimum caption duration (seconds)")
	flag.Float64Var(&captionMaxDuration, "maxcap", 7.0, "maximum caption duration (seconds)")
	flag.IntVar(&totalBuildsToKeep, "buildstokeep", 5, "total completed builds to keep when cleaning up")

	// just to clarify, by default this WILL remove files but adding the --keepoutput flag will cancel cleanup
	flag.BoolVar(&keepOutputFiles, "keepoutput", false, "keep artifacts after successful build")
	flag.BoolVar(&skipUpload, "skipupload", false, "skip upload step for builds")

	flag.Parse()

	if buildId == "" {
		buildId = util.Timestamp(time.Now())
	}

	logger.SetDefault(&logger.Config{
		Type:     logger.LoggerType(loggerType),
		Verbose:  verbose,
		WithAttr: []any{"buildId", buildId},
	})

	err := godotenv.Load()
	if err != nil {
		slog.Warn("error loading .env", "error", err)
	}
}

func main() {
	var err error

	ctx, cancel := context.WithCancel(context.Background())

	eh := errorhandler.DefaultErrorHandler(ctx, downloadCountPerPlaylist*3)
	defer eh.Report()

	// capture sigint
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		eh.Report()
		os.Exit(1)
	}()

	throttleInterval, err := time.ParseDuration(throttleDownloadsBy)
	if err != nil {
		panic(err)
	}

	// load config
	config := builder.Config{
		PlaylistIDs:                 strings.Split(playlistIDs, ","),
		ObjectStoreContainerName:    objectStoreContainerName,
		OutputDir:                   outputDir,
		MaxVideoSize:                uint(maxVideoSize),
		DownloadCountPerPlaylist:    downloadCountPerPlaylist,
		MaxPlaylistRequests:         maxPlaylistRequests,
		MinClipLengthSeconds:        minClipLengthSeconds,
		MaxClipLengthSeconds:        maxClipLengthSeconds,
		CaptionPrefixLength:         captionPrefixLength,
		CaptionNewsCorpusWeight:     captionNewsCorpusWeight,
		CaptionHospitalCorpusWeight: captionHospitalCorpusWeight,
		CaptionMinDuration:          captionMinDuration,
		CaptionMaxDuration:          captionMaxDuration,
		TotalBuildsToKeep:           totalBuildsToKeep,
		KeepOutputFiles:             keepOutputFiles,
		SkipUpload:                  skipUpload,
		ThrottleDownloadsBy:         throttleInterval,

		HospitalCorpus: hospitalText,
		OverlayImage:   overlayImage,
	}
	// auto-load env vars
	err = configloader.Load(&config)
	if err != nil {
		panic(err)
	}

	// for testing setting this to 1
	config.DownloadCountPerPlaylist = downloadCountPerPlaylist

	handleBuildPhaseErr := func(err error) {
		if err != nil {
			panic(fmt.Errorf("error building %q: %s", buildPhase, err))
		}
	}

	slog.Info("running", "buildPhase", buildPhase)

	switch domain.BuildPhase(buildPhase) {
	case domain.BuildPhaseAll:
		err = buildAll(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhaseSetup:
		err = buildSetup(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhaseDownloadVideos:
		err = downloadVideos(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhaseCutVideos:
		err = cutVideos(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhaseUploadVideos:
		err = uploadVideos(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhaseExtractImages:
		err = extractImages(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhasePosterImage:
		err = posterImage(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhaseModel:
		err = buildModel(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhaseGenerateCombinedVideo:
		err = buildCombinedVideo(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhaseManifest:
		err = buildManifest(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhasePromote:
		err = buildPromote(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhaseCleanup:
		err = buildCleanup(ctx, &config, eh)
		handleBuildPhaseErr(err)

	default:
		panic(fmt.Errorf("invalid build phase: %s", buildPhase))
	}

	slog.Info("Done")
}

func buildAll(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)
	return b.Run(ctx)
}

func buildSetup(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)
	return b.Setup(ctx)
}

func downloadVideos(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)

	err := b.Setup(ctx)
	if err != nil {
		return err
	}

	yesterday := util.Yesterday()
	videosStream := b.DownloadVideos(ctx, yesterday)

	videos := streams.StringSlice(ctx, videosStream)
	slog.Info("finished downloading videos", "count", len(videos))
	return nil
}

func cutVideos(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)

	paths, err := filepath.Glob("dist/*.webm")
	if err != nil {
		return fmt.Errorf("error getting video paths: %w", err)
	}
	if len(paths) == 0 {
		return errors.New("no videos found")
	}

	videoPathsStream := streams.StringStreamThrottled(ctx, time.Millisecond, paths...)
	clipsStream := b.CutVideos(ctx, videoPathsStream)

	clips := streams.StringSlice(ctx, clipsStream)
	slog.Info("finished processing clips", "count", len(clips))

	return nil
}

func uploadVideos(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)

	clips, err := filepath.Glob("dist/clips/*.webm")
	if err != nil {
		return fmt.Errorf("error getting clip paths: %w", err)
	}
	if len(clips) == 0 {
		return errors.New("no clip found")
	}

	uploadDir := buildId
	clipsStream := streams.StringStreamThrottled(ctx, time.Millisecond, clips...)
	uploadsStream := b.UploadVideos(ctx, uploadDir, clipsStream)

	uploads := streams.StringSlice(ctx, uploadsStream)
	slog.Info("finished uploading clips", "count", len(uploads))

	return nil
}

func extractImages(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)

	clips, err := filepath.Glob("dist/clips/*.webm")
	if err != nil {
		return fmt.Errorf("error getting clip paths: %w", err)
	}
	if len(clips) == 0 {
		return errors.New("no clip found")
	}

	clipsStream := streams.StringStream(ctx, clips...)
	imagesStream := b.ExtractImages(ctx, clipsStream)

	images := streams.StringSlice(ctx, imagesStream)
	slog.Info("finished extracting images", "count", len(images))

	return nil
}

func posterImage(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)

	images, err := filepath.Glob("dist/clips/*.png")
	if err != nil {
		return fmt.Errorf("error getting image paths: %w", err)
	}
	if len(images) == 0 {
		return errors.New("no clip found")
	}

	uploadDir := buildId

	posterImage, err := b.PosterImage(ctx, uploadDir, images)
	if err != nil {
		return fmt.Errorf("error generating poster image: %w", err)
	}

	slog.Info("finished outputting image", "image", posterImage)

	return nil
}

func buildModel(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)

	data, err := os.ReadFile("app/builder/hospital.txt")
	if err != nil {
		return err
	}

	uploadDir := getBuildID(config)

	fileName, err := b.Model(ctx, uploadDir, []domain.Corpus{
		{
			Type:     domain.CorpusTypeVTT,
			FileGlob: filepath.Join(config.OutputDir, "*.vtt"),
			Weight:   config.CaptionNewsCorpusWeight,
		},
		{
			Type:    domain.CorpusTypeString,
			Content: string(data),
			Weight:  config.CaptionHospitalCorpusWeight,
		},
	})
	if err != nil {
		return err
	}

	data, err = os.ReadFile(filepath.Join(config.OutputDir, fileName))
	if err != nil {
		return err
	}

	chain := captionschain.New(config.CaptionPrefixLength)
	err = chain.Load(data)
	if err != nil {
		return err
	}

	// cg := captiongenerator.New(chain, &captiongenerator.Config{
	// 	MinCaptionLength: 5,
	// 	MaxCaptionLength: 15,
	// })

	// prev := ""
	// for range 20 {
	// 	prev = cg.Caption(prev)
	// 	fmt.Println(prev)
	// }

	return nil
}

func buildCombinedVideo(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)

	clipFiles, err := filepath.Glob(filepath.Join(config.OutputDir, domain.ClipsDirName, "*.webm"))
	if err != nil {
		return err
	}

	modelFile := path.Join(config.OutputDir, domain.ModelFileName)

	videoFile, subsFile, err := b.GenerateCombinedVideo(ctx, domain.ArchivePrefix, modelFile, clipFiles)
	if err != nil {
		return err
	}

	slog.Info("finished creating combined video", "videoFile", videoFile, "subsFile", subsFile)

	return nil
}

func buildManifest(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)

	clips, err := filepath.Glob("dist/clips/*.webm")
	if err != nil {
		return err
	}

	buildDate := time.Now()
	buildID := buildId

	manifest := &domain.Manifest{
		BuildDate:   buildDate,
		ContentDate: util.Yesterday(),
		ID:          buildID,
		Files: domain.ManifestFiles{
			ModelFile: domain.ModelFileName,
			Clips:     make([]string, len(clips)),
		},
	}

	for i, clip := range clips {
		manifest.Files.Clips[i] = strings.Replace(clip, "dist/", "", 1)
	}

	result, err := b.Manifest(ctx, buildID, manifest)
	if err != nil {
		return err
	}

	slog.Info("successfully generated manifest", "file", result)

	return nil
}

func buildPromote(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)

	manifest, err := getManifest(config)
	if err != nil {
		return err
	}

	uploadDir := manifest.ID

	err = b.Promote(ctx, uploadDir)
	if err != nil {
		return err
	}

	slog.Info("finished promoting clips", "dir", uploadDir)
	return nil
}

func buildCleanup(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)

	deleted, err := b.Cleanup(ctx, config.TotalBuildsToKeep)
	if err != nil {
		return err
	}
	slog.Info("deleted objects", "count", len(deleted))
	return nil
}

func getBuildID(config *builder.Config) string {
	manifest, err := getManifest(config)
	if err != nil || manifest.ID == "" {
		return buildId
	}
	return manifest.ID
}

func getManifest(config *builder.Config) (*domain.Manifest, error) {

	data, err := os.ReadFile(filepath.Join(config.OutputDir, domain.ManifestFileName))
	if err != nil {
		return nil, fmt.Errorf("cannot open manifest file: %s", err)
	}

	manifest := domain.Manifest{}
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal manifest data: %s", err)
	}

	return &manifest, nil
}
