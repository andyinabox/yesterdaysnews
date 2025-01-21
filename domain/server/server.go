package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"sync"
	"text/template"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/assetshandler"
)

type Config struct {
	ObjectStoreUrl        string `env:"YN_OBJECTSTORE_URL" required:"true"`
	CDNUrl                string `env:"YN_CDN_URL" required:"true"`
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
	cg       domain.CaptionGenerator
	srv      *http.Server
	buildID  string
	manifest *domain.Manifest
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
