package main

import (
	"context"
	"flag"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/domain/imageprocessor"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/streams"
)

var inputGlob, outputFile string
var verbose bool

func init() {
	flag.StringVar(&inputGlob, "i", "dist/clips/*.png", "input file glob")
	flag.StringVar(&outputFile, "o", "dist/average.png", "output file")
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

const imageWidth = 1280
const imageHeight = 720

func main() {

	ctx := context.Background()

	eh := errorhandler.DefaultErrorHandler(ctx, 30)
	defer eh.Report()

	ip := imageprocessor.New()

	imagePaths, err := filepath.Glob(inputGlob)
	if err != nil {
		log.Fatal(err)
	}

	imagePathsStream := streams.StringStream(ctx, imagePaths...)

	log.Infof("averaging %d images", len(imagePaths))

	outputFile, err = ip.AverageImagesStream(ctx, eh.Channel(), imagePathsStream, domain.VideoWidth, domain.VideoHeight, outputFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("finished outputting %q", outputFile)

}
