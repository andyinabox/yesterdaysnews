package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/builder"
	"gitlab.com/andyinabox/yesterdaysnews/domain/captiongenerator"
	"gitlab.com/andyinabox/yesterdaysnews/domain/captionschain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/configloader"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

var verbose bool
var downloadCountPerPlaylist int
var buildPhase string

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&buildPhase, "b", "all", "build phase to execute")

	flag.IntVar(&downloadCountPerPlaylist, "d", 10, "download count per playlist")
	flag.Parse()

	// log.SetReportCaller(true)
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

	ctx := context.Background()

	eh := errorhandler.DefaultErrorHandler(ctx, downloadCountPerPlaylist*3)
	defer eh.Report()

	// load config
	config := builder.Config{}
	err := configloader.LoadJSONFile(&config, "builder.config.json")
	if err != nil {
		log.Fatal(err)
	}

	// for testing setting this to 1
	config.DownloadCountPerPlaylist = downloadCountPerPlaylist

	handleBuildPhaseErr := func(err error) {
		if err != nil {
			log.Fatalf("error building %s: %s", buildPhase, err)
		}
	}

	log.Info("running %q", buildPhase)

	switch domain.BuildPhase(buildPhase) {
	case domain.BuildPhaseAll:
		err = buildAll(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhaseSetup:
		err = buildSetup(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhaseVideoClips:
		err = buildVideoClips(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhaseModel:
		err = buildModel(ctx, &config, eh)
		handleBuildPhaseErr(err)
	case domain.BuildPhaseManifest:
		log.Fatalf("manifest build step not implemented")

	case domain.BuildPhasePromote:
		err = buildPromote(ctx, &config, eh)
		handleBuildPhaseErr(err)

	case domain.BuildPhaseCleanup:
		err = buildCleanup(ctx, &config, eh)
		handleBuildPhaseErr(err)

	default:
		log.Fatalf("invalid build phase: %s", buildPhase)
	}

	log.Info("Done")
}

func buildAll(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)
	return b.Run(ctx)
}

func buildSetup(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)
	return b.Setup(ctx)
}

func buildVideoClips(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)

	yesterday := util.Yesterday()
	uploadDir := util.Timestamp(yesterday)
	clips, err := b.VideoClips(ctx, yesterday, uploadDir)
	if err != nil {
		return err
	}
	log.Infof("finished processing %d clips", len(clips))
	return nil
}

func buildModel(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)

	uploadDir := getBuildID(config)

	fileName, err := b.Model(ctx, uploadDir, []domain.Corpus{
		{
			Type:     domain.CorpusTypeVTT,
			FileGlob: filepath.Join(config.OutputDir, "*.vtt"),
			Weight:   config.CaptionNewsCorpusWeight,
		},
		{
			Type:     domain.CorpusTypeText,
			FileGlob: "hospital.txt",
			Weight:   config.CaptionHospitalCorpusWeight,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(config.OutputDir, fileName))
	if err != nil {
		log.Fatal(err)
	}

	chain := captionschain.New(config.CaptionPrefixLength)
	err = chain.Load(data)
	if err != nil {
		log.Fatal(err)
	}

	cg := captiongenerator.New(chain, &captiongenerator.Config{
		MinCaptionLength: 5,
		MaxCaptionLength: 15,
	})

	prev := ""
	for i := 0; i < 20; i++ {
		prev = cg.Caption(prev)
		fmt.Println(prev)
	}

	return nil
}

func buildPromote(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)

	uploadDir := getBuildID(config)
	err := b.Promote(ctx, uploadDir)
	if err != nil {
		return err
	}

	log.Infof("finished promoting %q clips", uploadDir)
	return nil
}

func buildCleanup(ctx context.Context, config *builder.Config, eh domain.ErrorHandler) error {
	b := builder.New(config, eh)

	deleted, err := b.Cleanup(ctx, config.TotalBuildsToKeep)
	if err != nil {
		return err
	}
	log.Infof("deleted %d objects", len(deleted))
	return nil
}

func getBuildID(config *builder.Config) string {
	manifest, err := getManifest(config)
	if err != nil || manifest.ID == "" {
		return util.Timestamp(time.Now())
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
