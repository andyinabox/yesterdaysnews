package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"text/template"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/captionschain"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/textprocessor"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/uploader"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/assetshandler"
)

type Config struct {
	ObjectStoreUrl   string
	Templates        *template.Template
	Assets           fs.FS
	Port             int
	MinCaptionDelay  float64
	MaxCaptionDelay  float64
	MinCaptionLength int
	MaxCaptionLength int
}

type Server struct {
	tp       *textprocessor.Processor
	srv      *http.Server
	manifest *uploader.Manifest
	cfg      *Config
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
	handler.AddRoute("/clips", s.Clips())
	handler.AddRoute("/", s.Index())

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler,
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {
	var err error

	s.manifest, err = s.getManifest(ctx)
	if err != nil {
		return err
	}

	model, err := s.getModel(ctx)
	if err != nil {
		return err
	}

	chain := captionschain.New(0)
	err = chain.Load(model)
	if err != nil {
		return err
	}

	s.tp = textprocessor.New(chain, &textprocessor.Config{
		MinCaptionLength: s.cfg.MinCaptionLength,
		MaxCaptionLength: s.cfg.MaxCaptionLength,
	})

	log.Infof("starting server at http://localhost:%d", s.cfg.Port)
	return s.srv.ListenAndServe()
}

func (s *Server) getModel(ctx context.Context) ([]byte, error) {
	url := s.cfg.ObjectStoreUrl + "/" + s.manifest.Files.ModelFile

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (s *Server) getManifest(ctx context.Context) (*uploader.Manifest, error) {
	url := s.cfg.ObjectStoreUrl + "/manifest.json"

	log.Debug("create manifest request: " + url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	log.Debug("do manifest request")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	log.Debug("read manifest response")
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debug(string(data))

	log.Debug("unmarshal manifest data")
	manifest := uploader.Manifest{}
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}
