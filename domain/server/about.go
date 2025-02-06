package server

import (
	"html/template"
	"net/http"

	"github.com/charmbracelet/log"
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
			log.Error(err)
		}
	}

}
