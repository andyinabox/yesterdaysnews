package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/uploader"
)

func main() {

	outputDir := "dist"

	u := uploader.New(&uploader.Config{
		VideoFileName:    "yesterdays-news.mp4",
		SubsFileName:     "yesterdays-news.srt",
		CombinedFileName: "yesterdays-news-cc.mp4",
		ModelFileName:    "yesterdays-news.model.json",
		ClipsDirName:     "clips",
	})

	manifest, err := u.CreateManifest(context.Background(), outputDir)
	if err != nil {
		log.Fatal(err)
	}

	outFile := filepath.Join(outputDir, "manifest.json")

	b, err := json.Marshal(manifest)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("outputting manifest file", "file", outFile)
	err = os.WriteFile(outFile, b, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
