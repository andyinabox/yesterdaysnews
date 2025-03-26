package server

import (
	"html/template"
	"log/slog"
	"net/http"
)

type AboutRenderContext struct {
	RenderContext
	AboutContent template.HTML
}

func (s *Server) About() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data := AboutRenderContext{
			RenderContext: s.renderContext(),
			AboutContent:  template.HTML(s.cfg.AboutContent),
		}

		err := s.cfg.Templates.ExecuteTemplate(w, "about.html.tmpl", data)
		if err != nil {
			slog.Error("error rendering about page", "error", err)
		}
	}

}
