package main

import (
	"context"
	"flag"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/youtubedownloader"
)

var verbose bool
var listFilePath, outputDir string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&listFilePath, "f", "", "path to file with list of video ids")
	flag.StringVar(&outputDir, "o", "output/downloads", "where to download files")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	if listFilePath == "" {
		log.Fatal("no list file path provided")
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	dl := youtubedownloader.New(os.Getenv("YT_DLP_PATH"))

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(listFilePath)
	if err != nil {
		log.Fatal(err)
	}

	ids := strings.Split(strings.TrimSpace(string(data)), "\n")

	log.Debug(ids)

	log.Info("download videos", "listFilePath", listFilePath, "outputDir", outputDir)

	err = dl.DownloadVideoListWithDefaults(context.Background(), ids, outputDir)
	if err != nil {
		log.Fatal(err)
	}
}
