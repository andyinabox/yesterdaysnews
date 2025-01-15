package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/builder"
	"gitlab.com/andyinabox/yesterdaysnews/domain/captiongenerator"
	"gitlab.com/andyinabox/yesterdaysnews/domain/captionschain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/configloader"
)

var verbose bool

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}

	err := godotenv.Load()
	if err != nil {
		log.Warnf("error loading .env: %s", err)
	}
}

func main() {
	var b domain.Builder
	var eh domain.ErrorHandler

	ctx := context.Background()
	eh = errorhandler.DefaultErrorHandler(ctx, 0)
	defer eh.Report()

	// load config
	config := builder.Config{}
	err := configloader.LoadJSONFile(&config, "builder.config.json")
	if err != nil {
		log.Fatal(err)
	}

	b = builder.New(&config, eh)
	err = b.Setup(ctx)
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(config.OutputDir, domain.ManifestFileName))
	if err != nil {
		log.Fatalf("cannot open manifest file: %s", err)
	}

	manifest := domain.Manifest{}
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		log.Fatalf("cannot unmarshal manifest data: %s", err)
	}

	uploadDir := manifest.ID

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

	data, err = os.ReadFile(filepath.Join(config.OutputDir, fileName))
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
}
