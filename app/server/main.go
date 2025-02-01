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

//go:embed importmap.json
var importMap string

var verbose, loadAssetsFromFs bool
var port, prefixLength, minCaptionLength, maxCaptionLength int
var maxCaptionDelay, minCaptionDelay float64
var manifestCheckIntervalStr string

const assetsDirPath = "app/server/assets"
const assetsBuildDir = ".assets"

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

	ctx := context.Background()

	// strip out the name of the assets dir from the filesystem
	assetsEmbeddedFs, err := fs.Sub(fs.FS(assets), assetsBuildDir)
	if err != nil {
		log.Fatal(err)
	}

	// parse manifest check interval string into time.Duration
	manifestCheckInterval, err := time.ParseDuration(manifestCheckIntervalStr)
	if err != nil {
		log.Fatalf("error parsing manifest interval %s: %s", manifestCheckIntervalStr, err)
	}

	// parse about markdown
	aboutContent := blackfriday.Run(aboutContentMarkdown, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.Footnotes))

	cfg := &server.Config{
		Templates:             template.Must(template.ParseFS(templates, "tmpl/*")),
		AboutContent:          string(aboutContent),
		ImportMap:             importMap,
		AssetsDirFS:           os.DirFS(assetsDirPath),
		AssetsEmbeddedFS:      assetsEmbeddedFs,
		UseFilesystemAssets:   loadAssetsFromFs,
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
