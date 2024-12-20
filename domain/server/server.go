package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"text/template"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/assetshandler"
)

type CaptionGenerator interface {
	Caption(string) string
	MinCaptionLength() int
	MaxCaptionLength() int
}

type Config struct {
	ObjectStoreUrl  string
	Templates       *template.Template
	Assets          fs.FS
	Port            int
	MinCaptionDelay float64
	MaxCaptionDelay float64
}

type Server struct {
	cg  CaptionGenerator
	srv *http.Server
	cfg *Config
}

func New(cg CaptionGenerator, cfg *Config) *Server {

	s := &Server{
		cg:  cg,
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
	handler.AddRoute("/", s.Index())

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler,
	}

	return s
}

func (s *Server) ListenAndServe() error {
	log.Infof("starting server at http://localhost:%d", s.cfg.Port)
	return s.srv.ListenAndServe()
}
