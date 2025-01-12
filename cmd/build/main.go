package main

import (
	"context"
	"flag"

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
	var b domain.Builder
	var eh domain.ErrorHandler

	ctx := context.Background()

	eh = errorhandler.DefaultErrorHandler(ctx)
	defer eh.Report()

	// load config
	config := &builder.Config{}
	err := config.Load("builder.config.json")
	if err != nil {
		log.Fatal(err)
	}

	// for testing setting this to 1
	config.DownloadCountPerPlaylist = 1

	b = builder.New(config, eh)
	err = b.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Done")
}
