package main

import (
	"context"
	"flag"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/domain/videoprocessor"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/streams"
)

var clipsDir string
var verbose bool

func init() {
	flag.StringVar(&clipsDir, "d", "dist/clips", "dir to search for video clips")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	err := godotenv.Load()
	if err != nil {
		log.Warnf("error getting .env: %s", err)
	}
}

func main() {
	ctx := context.Background()

	eh := errorhandler.DefaultErrorHandler(ctx, 30)
	defer eh.Report()

	vp := videoprocessor.New(&videoprocessor.Config{})

	clips, err := filepath.Glob(filepath.Join(clipsDir, "*.webm"))
	if err != nil {
		log.Fatal(err)
	}

	clipsStream := streams.StringStream(ctx, clips...)

	imageFileStream := vp.ExtractImagesStream(ctx, eh.Channel(), clipsStream, time.Duration(0))

	imageFiles := streams.StringSlice(ctx, imageFileStream)

	log.Info("finished processing %d images", len(imageFiles))

}
