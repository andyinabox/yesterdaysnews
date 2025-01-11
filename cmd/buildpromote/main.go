package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/builder"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
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
	eh = errorhandler.DefaultErrorHandler(ctx)
	defer eh.DeferredReport()

	// load config
	config := &builder.Config{}
	err := config.Load("builder.config.json")
	if err != nil {
		log.Fatal(err)
	}

	b = builder.New(config, eh)
	err = b.Setup(ctx)
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(config.OutputDir, "manifest.json"))
	if err != nil {
		log.Fatalf("cannot open manifest file: %s", err)
	}

	manifest := domain.Manifest{}
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		log.Fatalf("cannot unmarshal manifest data: %s", err)
	}

	uploadDir := manifest.ID

	clips, err := b.Promote(ctx, uploadDir)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("finished processing %d clips", len(clips))
}
