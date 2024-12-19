package server

import (
	"html/template"
	"net/http"
)

type IndexRenderContext struct {
	JSVars template.JS
}

func (s *Server) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		jsVars := JSVars{
			ObjectStoreUrl: s.cfg.ObjectStoreUrl,
		}

		data := IndexRenderContext{
			JSVars: jsVars.ToJS(),
		}

		s.cfg.Templates.ExecuteTemplate(w, "index.html.tmpl", data)
	}
}
