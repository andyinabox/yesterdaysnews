package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/assetshandler"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

type Config struct {
	ObjectStoreUrl string `env:"YN_OBJECTSTORE_URL" required:"true"`
	CDNUrl         string `env:"YN_CDN_URL" required:"true"`
	Templates      *template.Template
	AboutContent   string
	ImportMap      string

	// see pkg/assetshandler for these options
	AssetsDirFS         fs.FS
	AssetsEmbeddedFS    fs.FS
	UseFilesystemAssets bool

	Port                  int
	MinCaptionDelay       float64
	MaxCaptionDelay       float64
	MinCaptionLength      int
	MaxCaptionLength      int
	ManifestCheckInterval time.Duration
}

type Server struct {
	cg         domain.CaptionGenerator
	srv        *http.Server
	buildID    string
	assetsPath string
	manifest   *domain.Manifest
	cfg        *Config
	reload     <-chan struct{}
	mu         sync.Mutex
}

func New(cfg *Config) *Server {

	s := &Server{
		cfg:        cfg,
		assetsPath: "/assets/" + util.TimestampBase64(time.Now(), false),
	}

	handler := assetshandler.New(
		&assetshandler.Config{
			AssetsRequestPathPrefix: s.assetsPath,
			AssetsEmbeddedFS:        cfg.AssetsEmbeddedFS,
			AssetsDirFS:             cfg.AssetsDirFS,
			UseFilesystemAssets:     cfg.UseFilesystemAssets,
		},
	)

	handler.AddRoute("/captions", s.Captions())
	handler.AddRoute("/reload", s.Reload())
	handler.AddRoute("/clips", s.Clips())
	handler.AddRoute("/about", s.About())
	handler.AddRoute("/", s.Index())

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler,
	}

	return s
}
