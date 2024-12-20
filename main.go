package main

import (
	"embed"
	"flag"
	"io/fs"
	"os"
	"text/template"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/captionschain"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/server"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/textprocessor"
)

//go:embed tmpl/*
var templates embed.FS

//go:embed assets/*
var assets embed.FS

var verbose, loadAssetsFromFs bool
var port, prefixLength, minCaptionLength, maxCaptionLength int
var maxCaptionDelay, minCaptionDelay float64
var objectStoreUrl string

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.BoolVar(&loadAssetsFromFs, "a", false, "load assets from filesystem (for easier frontend development)")
	flag.IntVar(&prefixLength, "p", 2, "markov chain prefix length")
	flag.IntVar(&minCaptionLength, "minl", 5, "min caption length in words")
	flag.IntVar(&maxCaptionLength, "maxl", 15, "max caption length in words")
	flag.Float64Var(&minCaptionDelay, "mind", 1.5, "min caption delay in seconds")
	flag.Float64Var(&maxCaptionDelay, "maxd", 5.0, "max caption delay in seconds")
	flag.IntVar(&port, "port", 8080, "server port")
	flag.StringVar(&objectStoreUrl, "url", "https://localhost:9000", "url of object storage")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
		log.SetReportTimestamp(false)
	}

}

func main() {
	var assetsFs fs.FS

	// load assets from filesystem for development
	if loadAssetsFromFs {
		assetsFs = os.DirFS("assets")
		// load assets from embedded data
	} else {
		var err error
		assetsFs, err = fs.Sub(fs.FS(assets), "assets")
		if err != nil {
			log.Fatal(err)
		}
	}

	data, err := os.ReadFile("dist/yesterdays-news.model.json")
	if err != nil {
		log.Fatal(err)
	}

	chain := captionschain.New(prefixLength)
	err = chain.Load(data)
	if err != nil {
		log.Fatal(err)
	}

	tp := textprocessor.New(chain, &textprocessor.Config{
		MinCaptionLength: minCaptionLength,
		MaxCaptionLength: maxCaptionLength,
	})

	s := server.New(tp, &server.Config{
		ObjectStoreUrl:  objectStoreUrl,
		Templates:       template.Must(template.New("").ParseFS(templates, "tmpl/*.tmpl")),
		Assets:          assetsFs,
		Port:            port,
		MinCaptionDelay: minCaptionDelay,
		MaxCaptionDelay: maxCaptionDelay,
	})

	log.Fatal(s.ListenAndServe())
}
