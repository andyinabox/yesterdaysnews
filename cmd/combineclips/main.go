package main

import (
	"context"
	"flag"
	"path/filepath"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/videoprocessor"
)

var verbose bool
var inputGlob, outputFile string

func init() {

	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&inputGlob, "i", "", "glob to get input videos")
	flag.StringVar(&outputFile, "o", "", "name of output file (path/to/video.mp4)")

	flag.Parse()

	if inputGlob == "" || outputFile == "" {
		log.Fatal("must set input and output")
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}
}
func main() {

	var vp domain.VideoProcessor

	vp = videoprocessor.New(&videoprocessor.Config{})

	files, err := filepath.Glob(inputGlob)
	if err != nil {
		log.Fatal(err)
	}

	file, err := vp.ShuffleClipsAndCombine(context.Background(), files, outputFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("finished combining clips into %s", file)

}
