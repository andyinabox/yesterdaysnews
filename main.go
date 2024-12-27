package main

import (
	"context"
	"embed"
	"flag"
	"io/fs"
	"os"
	"text/template"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/server"
)

//go:embed tmpl/*
var templates embed.FS

//go:embed assets/*
var assets embed.FS

var verbose, loadAssetsFromFs bool
var port, prefixLength, minCaptionLength, maxCaptionLength int
var maxCaptionDelay, minCaptionDelay float64
var manifestCheckIntervalStr string

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.BoolVar(&loadAssetsFromFs, "a", false, "load assets from filesystem (for easier frontend development)")
	// flag.BoolVar(&loadObjectsFromFs, "o", true, "load objects from filesystem")
	flag.IntVar(&prefixLength, "p", 2, "markov chain prefix length")
	flag.IntVar(&minCaptionLength, "minl", 5, "min caption length in words")
	flag.IntVar(&maxCaptionLength, "maxl", 15, "max caption length in words")
	flag.Float64Var(&minCaptionDelay, "mind", 1.5, "min caption delay in seconds")
	flag.Float64Var(&maxCaptionDelay, "maxd", 5.0, "max caption delay in seconds")
	flag.IntVar(&port, "port", 8080, "server port")
	flag.StringVar(&manifestCheckIntervalStr, "m", "1h", "manifest check interval")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
		log.SetReportTimestamp(false)
		log.SetReportCaller(true)
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

	manifestCheckInterval, err := time.ParseDuration(manifestCheckIntervalStr)
	if err != nil {
		log.Fatalf("error paring manifest interval %s: %s", manifestCheckIntervalStr, err)
	}

	cfg := &server.Config{
		ObjectStoreUrl:        os.Getenv("YN_OBJECTSTORE_URL"),
		Templates:             template.Must(template.New("").ParseFS(templates, "tmpl/*.tmpl")),
		Assets:                assetsFs,
		Port:                  port,
		MinCaptionDelay:       minCaptionDelay,
		MaxCaptionDelay:       maxCaptionDelay,
		MinCaptionLength:      minCaptionLength,
		MaxCaptionLength:      maxCaptionLength,
		ManifestCheckInterval: manifestCheckInterval,
	}

	log.Info("creating new server", "config", cfg)

	s := server.New(cfg)

	log.Fatal(s.Start(context.Background()))
}
