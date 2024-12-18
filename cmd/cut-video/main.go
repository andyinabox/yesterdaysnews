package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"

	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/videoprocessor"
)

var verbose, outputManifest bool
var inputPath, outputDir string
var minLength, maxLength int

func init() {

	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.BoolVar(&outputManifest, "m", true, "output manifest json")
	flag.StringVar(&inputPath, "i", "", "path to input video")
	flag.StringVar(&outputDir, "o", "", "path to output dir")
	flag.IntVar(&minLength, "min", 5, "minimum clip length")
	flag.IntVar(&maxLength, "max", 20, "maximum clip length")
	flag.Parse()

	if inputPath == "" || outputDir == "" {
		log.Fatal("must set input and output paths")
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}
}

func main() {

	vp := videoprocessor.New(&videoprocessor.Config{})

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	result, err := vp.BreakVideoIntoClips(context.Background(), inputPath, outputDir, minLength, maxLength)
	if err != nil {
		log.Fatal(err)
	}

	if outputManifest {
		outFile := filepath.Join(outputDir, "manifest.json")

		b, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("outputting manifest file", "file", outFile)
		err = os.WriteFile(outFile, b, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Infof("finished generating %d videos", len(result.Files))

}
