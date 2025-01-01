package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/videoprocessor"
)

var verbose, outputManifest bool
var inputGlob, outputDir string
var minLength, maxLength int

func init() {

	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.BoolVar(&outputManifest, "m", true, "output manifest json")
	flag.StringVar(&inputGlob, "i", "", "glob to get input videos")
	flag.StringVar(&outputDir, "o", "", "path to output dir")
	flag.IntVar(&minLength, "min", 5, "minimum clip length")
	flag.IntVar(&maxLength, "max", 20, "maximum clip length")
	flag.Parse()

	if inputGlob == "" || outputDir == "" {
		log.Fatal("must set input and output paths")
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}
func main() {

	var vp domain.VideoProcessor

	vp = videoprocessor.New(&videoprocessor.Config{})

	files, err := filepath.Glob(inputGlob)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	result := domain.BreakVideoIntoClipsResult{
		Files: []string{},
	}

	for _, f := range files {
		log.Infof("cutting up video %s", f)

		r, err := vp.BreakVideoIntoClips(context.Background(), f, outputDir, minLength, maxLength)
		if err != nil {
			log.Error(err)
			continue
		}

		result.Files = append(result.Files, r.Files...)
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

	log.Infof("finished outputting %d videos", len(result.Files))

}
