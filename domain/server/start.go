package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/captiongenerator"
	"gitlab.com/andyinabox/yesterdaysnews/domain/captionschain"
)

func (s *Server) Start(ctx context.Context) error {
	var err error

	ctx, cancel := context.WithCancel(ctx)
	reload := make(chan struct{})
	s.reload = reload

	defer func() {
		cancel()
		close(reload)
	}()

	s.buildID, err = s.getCurrentBuildID(ctx)
	if err != nil {
		return fmt.Errorf("error loading buildID on startup: %w", err)
	}

	s.manifest, err = s.getManifest(ctx)
	if err != nil {
		return fmt.Errorf("error loading manifest on startup: %w", err)
	}

	s.cg, err = s.makeCaptionGenerator(ctx)
	if err != nil {
		return fmt.Errorf("error making caption generator on startup: %w", err)
	}

	// poll for changes to manifest
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(s.cfg.ManifestCheckInterval)
				log.Debug("checking manifest")

				buildID, err := s.getCurrentBuildID(ctx)
				if err != nil {
					log.Errorf("error loading current buildID: %s", err)
					continue
				}

				log.Debug("comparing buildID", "current", s.buildID, "new", buildID)
				if buildID != s.buildID {
					log.Info("buildID is updated, reloading...")

					manifest, err := s.getManifest(ctx)
					if err != nil {
						log.Errorf("error loading manifest: %s", err)
						continue
					}

					s.mu.Lock()
					s.buildID = buildID
					s.manifest = manifest
					s.mu.Unlock()

					cg, err := s.makeCaptionGenerator(ctx)
					if err != nil {
						log.Errorf("error making new text processor: %s", err)
						continue
					}

					s.mu.Lock()
					s.cg = cg
					s.mu.Unlock()

					reload <- struct{}{}
				}
			}
		}
	}()

	log.Infof("starting server at http://localhost:%d", s.cfg.Port)
	return s.srv.ListenAndServe()
}

func (s *Server) makeCaptionGenerator(ctx context.Context) (domain.CaptionGenerator, error) {
	model, err := s.getModel(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting model for caption generator: %w", err)
	}

	chain := captionschain.New(0)
	err = chain.Load(model)
	if err != nil {
		return nil, fmt.Errorf("error loading model for caption generator: %w", err)
	}

	cg := captiongenerator.New(chain, &captiongenerator.Config{
		MinCaptionLength: s.cfg.MinCaptionLength,
		MaxCaptionLength: s.cfg.MaxCaptionLength,
	})

	return cg, nil
}

func (s *Server) getModel(ctx context.Context) ([]byte, error) {
	// url := s.cfg.ObjectStoreUrl + "/" + s.manifest.Files.ModelFile
	url := fmt.Sprintf("%s/%s/%s", s.cfg.ObjectStoreUrl, s.buildID, s.manifest.Files.ModelFile)

	log.Debug("create model request: " + url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("recieved non-200 status code: %q", resp.Status)
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (s *Server) getManifest(ctx context.Context) (*domain.Manifest, error) {

	url := fmt.Sprintf("%s/%s/manifest.json", s.cfg.ObjectStoreUrl, s.buildID)

	log.Debug("create manifest request: " + url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// log.Debug("do manifest request")
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
	manifest := domain.Manifest{}
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (s *Server) getCurrentBuildID(ctx context.Context) (string, error) {
	url := s.cfg.ObjectStoreUrl + "/current.txt"

	log.Debug("create buildID request: " + url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	log.Debug("do buildID request")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("recieved non-200 status code: %q", resp.Status)
	}

	log.Debug("read buildID response")
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Debug(string(data))

	return string(data), nil
}
