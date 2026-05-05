package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
	"gitlab.com/andyinabox/yesterdaysnews/domain/logger"
	"gitlab.com/andyinabox/yesterdaysnews/domain/server"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/configloader"
)

//go:embed tmpl/*.tmpl
var templates embed.FS

//go:embed about.md
var aboutContentMarkdown []byte

//go:embed description.txt
var siteDescription string

//go:embed .assets/*
var assets embed.FS

//go:embed importmap.json
var importMap string

var verbose, loadAssetsFromFs bool
var loggerType string
var port, prefixLength, minCaptionLength, maxCaptionLength, fetchClipsWhenLowerThan int
var maxCaptionDelay, minCaptionDelay float64
var manifestCheckIntervalStr string

const assetsDirPath = "app/server/assets"
const assetsBuildDir = ".assets"

func init() {

	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.StringVar(&loggerType, "log", "text", "logger type (text, json)")
	flag.BoolVar(&loadAssetsFromFs, "a", false, "load assets from filesystem (for easier frontend development)")
	flag.IntVar(&prefixLength, "p", 2, "markov chain prefix length")
	flag.IntVar(&minCaptionLength, "minl", 7, "min caption length in words")
	flag.IntVar(&maxCaptionLength, "maxl", 15, "max caption length in words")
	flag.IntVar(&fetchClipsWhenLowerThan, "fetchclips", 5, "on frontend, when clips count is lower than this fetch more")
	flag.Float64Var(&minCaptionDelay, "mind", 3, "min caption delay in seconds")
	flag.Float64Var(&maxCaptionDelay, "maxd", 7.0, "max caption delay in seconds")
	flag.IntVar(&port, "port", 8080, "server port")
	flag.StringVar(&manifestCheckIntervalStr, "m", "1h", "manifest check interval")
	flag.Parse()

	logger.SetDefault(&logger.Config{
		Type:    logger.LoggerType(loggerType),
		Verbose: verbose,
	})

}

func main() {

	ctx := context.Background()

	// strip out the name of the assets dir from the filesystem
	assetsEmbeddedFs, err := fs.Sub(fs.FS(assets), assetsBuildDir)
	if err != nil {
		panic(err)
	}

	// parse manifest check interval string into time.Duration
	manifestCheckInterval, err := time.ParseDuration(manifestCheckIntervalStr)
	if err != nil {
		panic(fmt.Sprintf("error parsing manifest interval %s: %s", manifestCheckIntervalStr, err))
	}

	// parse about markdown
	aboutContent := blackfriday.Run(aboutContentMarkdown, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.Footnotes))

	cfg := &server.Config{
		Templates:               template.Must(template.ParseFS(templates, "tmpl/*")),
		AboutContent:            string(aboutContent),
		ImportMap:               importMap,
		SiteDescription:         strings.TrimSpace(siteDescription),
		AssetsDirFS:             os.DirFS(assetsDirPath),
		AssetsEmbeddedFS:        assetsEmbeddedFs,
		UseFilesystemAssets:     loadAssetsFromFs,
		Port:                    port,
		MinCaptionDelay:         minCaptionDelay,
		MaxCaptionDelay:         maxCaptionDelay,
		MinCaptionLength:        minCaptionLength,
		MaxCaptionLength:        maxCaptionLength,
		ManifestCheckInterval:   manifestCheckInterval,
		FetchClipsWhenLowerThan: fetchClipsWhenLowerThan,
	}

	err = configloader.Load(cfg)
	if err != nil {
		panic(fmt.Sprintf("error loading env vars: %s", err))
	}

	slog.Info("creating new server")

	s := server.New(cfg)

	err = s.Start(ctx)
	if err != nil {
		panic(err)
	}
}
