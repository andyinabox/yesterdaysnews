package main

import (
	"flag"
	"math"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/imgavg"
)

var inputGlob, outputFile, workDir string
var maxImages int
var verbose bool

func init() {
	flag.StringVar(&inputGlob, "i", "dist/clips/*.png", "input file glob")
	flag.StringVar(&outputFile, "o", "dist/average-test.png", "output file")
	flag.StringVar(&workDir, "w", "dist/avg", "working dir for intermediate files")
	flag.IntVar(&maxImages, "max", math.MaxInt, "max images to process")
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

	err := os.MkdirAll(workDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	paths, err := filepath.Glob(inputGlob)
	if err != nil {
		log.Fatal(err)
	}

	client := imgavg.New(func(err error) {
		log.Error(err)
	})

	if maxImages > len(paths) {
		maxImages = len(paths)
	}

	result, err := client.Run(paths[:maxImages], outputFile, imageWidth, imageHeight)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("successfully outputed image %q", result)
}

// func main() {

// 	ctx := context.Background()

// 	eh := errorhandler.DefaultErrorHandler(ctx, 30)
// 	defer eh.Report()

// 	ip := imageprocessor.New()

// 	imagePaths, err := filepath.Glob(inputGlob)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	imagePathsStream := streams.StringStream(ctx, imagePaths...)

// 	log.Infof("averaging %d images", len(imagePaths))

// 	outputFile, err = ip.AverageImagesStream(ctx, eh.Channel(), imagePathsStream, domain.VideoWidth, domain.VideoHeight, outputFile)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Infof("finished outputting %q", outputFile)

// }
