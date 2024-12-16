package main

import (
	"context"
	"flag"
	"os"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/youtubedownloader"
)

var verbose bool
var url, outputDir string

func init() {
	flag.StringVar(&url, "u", "", "video url or id")
	flag.StringVar(&outputDir, "o", "output", "output dir")
	flag.BoolVar(&verbose, "v", true, "verbose output")
	flag.Parse()

	if url == "" {
		log.Fatal("no url provided")
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {

	dl := youtubedownloader.New(os.Getenv("YT_DLP_PATH"))

	err := dl.DownloadVideoWithDefaults(context.Background(), url, outputDir)
	if err != nil {
		log.Fatal(err)
	}

}
