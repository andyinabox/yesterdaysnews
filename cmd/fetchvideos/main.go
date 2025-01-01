package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/downloader"
)

var verbose, outputManifest bool
var channelName, outputDir string
var maxResults int

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&channelName, "n", "@CNN", "username/handle for channel")
	flag.StringVar(&outputDir, "o", "output/downloads/cnn", "where to download files")
	flag.IntVar(&maxResults, "c", 10, "max number of videos to download")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.BoolVar(&outputManifest, "m", true, "output manifest file")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}

}

func main() {

	var dl domain.Downloader

	dl = downloader.New(&downloader.Config{
		GoogleAPIKey: os.Getenv("GOOGLE_API_KEY"),
		BinPathYTDLP: os.Getenv("YT_DLP_PATH"),
	})

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	result, err := dl.DownloadVideosForChannel(context.Background(), domain.DownloadRequest{
		ChannelUsername: channelName,
		Date:            time.Now().AddDate(0, 0, -1),
		MaxResults:      maxResults,
		OutputDir:       outputDir,
	})
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

	log.Infof("finished downloading %d videos", len(result.Files))
}
