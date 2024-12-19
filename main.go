package main

import (
	"embed"
	"flag"
	"io/fs"
	"os"
	"text/template"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/server"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/textprocessor"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/markov"
)

//go:embed tmpl/*
var templates embed.FS

//go:embed assets/*
var assets embed.FS

var verbose, loadAssetsFromFs bool
var port, prefixLength, maxCaptionDelay, minCaptionDelay int
var objectStoreUrl string

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.BoolVar(&loadAssetsFromFs, "a", false, "load assets from filesystem (for easier frontend development)")
	flag.IntVar(&prefixLength, "p", 2, "markov chain prefix length")
	flag.IntVar(&minCaptionDelay, "min", 1, "min caption delay")
	flag.IntVar(&maxCaptionDelay, "max", 6, "max caption delay")
	flag.IntVar(&port, "port", 8080, "server port")
	flag.StringVar(&objectStoreUrl, "o", "https://localhost:9000", "url of object storage")
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

	file, err := os.Open("data/hospital.txt")
	if err != nil {
		log.Fatal(err)
	}

	chain := markov.NewBasicChain(prefixLength)
	chain.Build(file)

	tp := textprocessor.New(chain, &textprocessor.Config{
		// PrefixLength:     prefixLength,
		MinCaptionLength: 5,
		MaxCaptionLength: 10,
	})

	s := server.New(tp, &server.Config{
		ObjectStoreUrl:  objectStoreUrl,
		Templates:       template.Must(template.New("").ParseFS(templates, "tmpl/*.tmpl")),
		Assets:          assetsFs,
		Port:            port,
		MinCaptionDelay: time.Duration(minCaptionDelay) * time.Second,
		MaxCaptionDelay: time.Duration(maxCaptionDelay) * time.Second,
	})

	log.Fatal(s.ListenAndServe())
}
