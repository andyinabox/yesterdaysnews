package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"sync"
	"text/template"
	"time"

	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/manifest"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/textprocessor"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/assetshandler"
)

type Config struct {
	ObjectStoreUrl        string
	Templates             *template.Template
	Assets                fs.FS
	Port                  int
	MinCaptionDelay       float64
	MaxCaptionDelay       float64
	MinCaptionLength      int
	MaxCaptionLength      int
	ManifestCheckInterval time.Duration
}

type Server struct {
	tp       *textprocessor.Processor
	srv      *http.Server
	manifest *manifest.Manifest
	cfg      *Config
	reload   <-chan struct{}
	mu       sync.Mutex
}

func New(cfg *Config) *Server {

	s := &Server{
		cfg: cfg,
	}

	handler := assetshandler.New(
		&assetshandler.Config{
			AssetsUrlPath:     "/assets",
			AssetsFS:          cfg.Assets,
			StripAssetsPrefix: true,
		},
	)

	handler.AddRoute("/captions", s.Captions())
	handler.AddRoute("/reload", s.Reload())
	handler.AddRoute("/clips", s.Clips())
	handler.AddRoute("/", s.Index())

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler,
	}

	return s
}
