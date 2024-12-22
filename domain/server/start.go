package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/captionschain"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/textprocessor"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/uploader"
)

const ManifestCheckInterval = 5 * time.Second

func (s *Server) Start(ctx context.Context) error {
	var err error

	ctx, cancel := context.WithCancel(ctx)
	reload := make(chan struct{})
	s.reload = reload

	defer func() {
		cancel()
		close(reload)
	}()

	s.manifest, err = s.getManifest(ctx)
	if err != nil {
		return err
	}

	s.tp, err = s.makeTextProcesor(ctx)
	if err != nil {
		return err
	}

	// poll for changes to manifest
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(ManifestCheckInterval)
				log.Debug("checking manifest")

				manifest, err := s.getManifest(ctx)
				if err != nil {
					log.Errorf("error loading manifest: %s", err)
					continue
				}

				// log.Debug("comparing manifest date", "current", s.manifest.Date, "new", manifest.Date)
				if manifest.Date != s.manifest.Date {
					log.Info("manifest is updated, reloading...")

					s.mu.Lock()
					s.manifest = manifest
					s.mu.Unlock()

					tp, err := s.makeTextProcesor(ctx)
					if err != nil {
						log.Errorf("error making new text processor: %s", err)
						continue
					}

					s.mu.Lock()
					s.tp = tp
					s.mu.Unlock()

					reload <- struct{}{}
				}
			}
		}
	}()

	log.Infof("starting server at http://localhost:%d", s.cfg.Port)
	return s.srv.ListenAndServe()
}

func (s *Server) makeTextProcesor(ctx context.Context) (*textprocessor.Processor, error) {
	model, err := s.getModel(ctx)
	if err != nil {
		return nil, err
	}

	chain := captionschain.New(0)
	err = chain.Load(model)
	if err != nil {
		return nil, err
	}

	tp := textprocessor.New(chain, &textprocessor.Config{
		MinCaptionLength: s.cfg.MinCaptionLength,
		MaxCaptionLength: s.cfg.MaxCaptionLength,
	})

	return tp, nil
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

	// log.Debug("create manifest request: " + url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// log.Debug("do manifest request")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// log.Debug("read manifest response")
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// log.Debug(string(data))

	// log.Debug("unmarshal manifest data")
	manifest := uploader.Manifest{}
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}
