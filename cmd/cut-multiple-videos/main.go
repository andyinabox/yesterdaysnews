package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/videoeditor"
)

var verbose bool
var inputGlob, outputDir string
var minLength, maxLength int

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	flag.BoolVar(&verbose, "v", false, "verbose output")
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
	files, err := filepath.Glob(inputGlob)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	editor := videoeditor.New(os.Getenv("FFMPEG_PATH"), os.Getenv("FFPROBE_PATH"))

	log.Debug("parse glob", "glob", inputGlob, "files", files)

	for _, f := range files {
		err = editor.BreakVideoIntoClips(context.Background(), f, outputDir, minLength, maxLength)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

}
