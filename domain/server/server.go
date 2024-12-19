package server

import (
	"embed"
	"fmt"
	"net/http"
	"text/template"
)

type TextGenerator interface {
	Sentence() string
}

type Config struct {
	ObjectStoreUrl string
	PublicFiles    embed.FS
	Templates      *template.Template
	Port           int
}

type Server struct {
	tg  TextGenerator
	srv *http.Server
	cfg *Config
}

func New(tg TextGenerator, cfg *Config) *Server {

	s := &Server{
		tg:  tg,
		cfg: cfg,
	}

	mux := http.NewServeMux()
	mux.Handle("/captions", s.Captions())
	mux.Handle("/", s.Index())

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	return s
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}
