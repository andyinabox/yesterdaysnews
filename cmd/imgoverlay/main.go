package main

import (
	"context"
	"flag"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/domain/imageprocessor"
)

var baseImagePath, overlayImagePath, outputImagePath string
var verbose bool

func init() {
	flag.StringVar(&baseImagePath, "base", "dist/average.png", "base image on bottom layer")
	flag.StringVar(&overlayImagePath, "overlay", "app/builder/overlay.png", "overlay image on top layer")
	flag.StringVar(&outputImagePath, "o", "dist/poster.png", "output image location")
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

	baseImgFile, err := os.Open(baseImagePath)
	if err != nil {
		log.Debug(err)
	}

	overlayImageFile, err := os.Open(overlayImagePath)
	if err != nil {
		log.Debug(err)
	}

	ctx := context.Background()

	eh := errorhandler.DefaultErrorHandler(ctx, 30)
	defer eh.Report()

	ip := imageprocessor.New()

	output, err := ip.OverlayImages(ctx, baseImgFile, overlayImageFile, outputImagePath)

	log.Infof("finished outputing %q", output)

}
