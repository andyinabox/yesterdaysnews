package main

import (
	"context"
	"flag"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/builder"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/configloader"
)

var verbose bool
var downloadCountPerPlaylist int

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose output")
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
	var b domain.Builder
	var eh domain.ErrorHandler

	ctx := context.Background()

	eh = errorhandler.DefaultErrorHandler(ctx, downloadCountPerPlaylist*3)
	defer eh.Report()

	// load config
	config := builder.Config{}
	err := configloader.LoadJSONFile(&config, "builder.config.json")
	if err != nil {
		log.Fatal(err)
	}

	// for testing setting this to 1
	config.DownloadCountPerPlaylist = downloadCountPerPlaylist

	b = builder.New(&config, eh)
	err = b.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Done")
}
