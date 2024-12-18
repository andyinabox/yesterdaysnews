package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/textprocessor"
)

var verbose bool
var glob string

func init() {
	flag.StringVar(&glob, "f", "", "glob for files to process")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	if glob == "" {
		log.Fatal("no glob")
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {

	tp := textprocessor.New(&textprocessor.Config{})

	files, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("parse glob", "glob", glob, "files", files)

	corpus := ""

	ctx := context.Background()

	for _, f := range files {
		text, err := tp.VTTToText(ctx, f)
		if err != nil {
			log.Error(err)
			continue
		}

		corpus = corpus + "\n" + text
	}

	fmt.Print(corpus)

}
