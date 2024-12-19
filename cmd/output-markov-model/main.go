package main

import (
	"flag"
	"os"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/textprocessor"
)

var verbose bool
var prefixLength int
var modelPath string

func init() {

	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.IntVar(&prefixLength, "p", 2, "prefix length")
	flag.StringVar(&modelPath, "m", "dist/yesterdays-news.model.json", "path to export the model file")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}
}

func main() {

	p := textprocessor.New(&textprocessor.Config{
		PrefixLength: prefixLength,
	})

	err := p.Build([]textprocessor.Corpus{
		{
			Type:     textprocessor.CorpusTypeVTT,
			FileGlob: "download/cnn/*.vtt",
			Weight:   1.0,
		},
		{
			Type:     textprocessor.CorpusTypeVTT,
			FileGlob: "download/msnbc/*.vtt",
			Weight:   1.0,
		},
		{
			Type:     textprocessor.CorpusTypeVTT,
			FileGlob: "download/foxnews/*.vtt",
			Weight:   1.0,
		},
		{
			Type:     textprocessor.CorpusTypeText,
			FileGlob: "data/hospital.txt",
			Weight:   1.0,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	modelJSON, err := p.ExportModel()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("outputting model file", "file", modelPath)
	err = os.WriteFile(modelPath, modelJSON, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

}
