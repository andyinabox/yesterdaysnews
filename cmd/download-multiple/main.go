package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/youtubedownloader"
)

var verbose, outputManifest bool
var listFilePath, outputDir string

type manifestEntry struct {
	Video string `json:"video"`
	Subs  string `json:"subs"`
}

type manifest struct {
	Files []manifestEntry `json:"files"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&listFilePath, "f", "", "path to file with list of video ids")
	flag.StringVar(&outputDir, "o", "output/downloads", "where to download files")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.BoolVar(&outputManifest, "m", true, "output manifest file")
	flag.Parse()

	if listFilePath == "" {
		log.Fatal("no list file path provided")
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
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

	// output manifest file
	if outputManifest {

		files, err := os.ReadDir(outputDir)
		if err != nil {
			log.Fatal("error reading output dir, cannot output manifest", "error", err)
		}

		m := manifest{
			Files: make([]manifestEntry, len(ids)),
		}

		// expected manifest
		for i, id := range ids {

			e := manifestEntry{
				Video: id + ".mp4",
				Subs:  id + ".en.vtt",
			}

			var videoFound bool
			var subsFound bool

			// check manifest against fs
			for _, f := range files {
				if f.Name() == e.Video {
					videoFound = true
				}
				if f.Name() == e.Subs {
					subsFound = true
				}
			}

			if !videoFound {
				log.Errorf("expected file %s not found in fs", e.Video)
			}

			if !subsFound {
				log.Errorf("expected file %s not found in fs", e.Subs)
			}

			m.Files[i] = e
		}

		outFile := filepath.Join(outputDir, "manifest.json")

		b, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("outputting manifest file", "file", outFile)
		err = os.WriteFile(outFile, b, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}
