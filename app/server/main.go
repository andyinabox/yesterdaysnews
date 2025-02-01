package main

import (
	"context"
	"embed"
	"flag"
	"html/template"
	"io/fs"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/russross/blackfriday/v2"
	"gitlab.com/andyinabox/yesterdaysnews/domain/server"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/configloader"
)

//go:embed tmpl/*.tmpl
var templates embed.FS

//go:embed about.md
var aboutContentMarkdown []byte

//go:embed .assets/*
var assets embed.FS

var verbose, loadAssetsFromFs bool
var port, prefixLength, minCaptionLength, maxCaptionLength int
var maxCaptionDelay, minCaptionDelay float64
var manifestCheckIntervalStr string

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Warnf("error loading .env file: %s", err)
	}

	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.BoolVar(&loadAssetsFromFs, "a", false, "load assets from filesystem (for easier frontend development)")
	flag.IntVar(&prefixLength, "p", 2, "markov chain prefix length")
	flag.IntVar(&minCaptionLength, "minl", 5, "min caption length in words")
	flag.IntVar(&maxCaptionLength, "maxl", 12, "max caption length in words")
	flag.Float64Var(&minCaptionDelay, "mind", 2, "min caption delay in seconds")
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

	ctx := context.Background()

	// load assets from filesystem for development
	if loadAssetsFromFs {
		assetsFs = os.DirFS("app/server/assets")
		// load assets from embedded data
	} else {
		var err error
		assetsFs, err = fs.Sub(fs.FS(assets), ".assets")
		if err != nil {
			log.Fatal(err)
		}
	}

	manifestCheckInterval, err := time.ParseDuration(manifestCheckIntervalStr)
	if err != nil {
		log.Fatalf("error parsing manifest interval %s: %s", manifestCheckIntervalStr, err)
	}

	aboutContent := blackfriday.Run(aboutContentMarkdown, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.Footnotes))

	cfg := &server.Config{
		Templates:             template.Must(template.ParseFS(templates, "tmpl/*")),
		AboutContent:          string(aboutContent),
		Assets:                assetsFs,
		Port:                  port,
		MinCaptionDelay:       minCaptionDelay,
		MaxCaptionDelay:       maxCaptionDelay,
		MinCaptionLength:      minCaptionLength,
		MaxCaptionLength:      maxCaptionLength,
		ManifestCheckInterval: manifestCheckInterval,
	}

	err = configloader.Load(cfg)
	if err != nil {
		log.Fatalf("error loading env vars: %s", err)
	}

	log.Info("creating new server")
	log.Infof("%#v", cfg)

	s := server.New(cfg)

	log.Fatal(s.Start(ctx))
}
