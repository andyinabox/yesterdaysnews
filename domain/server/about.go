package server

import (
	"html/template"
	"net/http"
)

type AboutRenderContext struct {
	PageTitle    string
	AboutContent template.HTML
}

func (s *Server) About() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		title := "yesterday's news"

		data := AboutRenderContext{
			PageTitle:    title,
			AboutContent: template.HTML(s.cfg.AboutContent),
		}

		s.cfg.Templates.ExecuteTemplate(w, "about.html.tmpl", data)
	}

}
