package server

import (
	"fmt"
	"math/rand"
	"net/http"
)

type IndexRenderContext struct {
	InitialClipURL string
	InitialCaption string
}

func (s *Server) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		clipUrl := fmt.Sprintf(
			"%s/%s",
			s.cfg.ObjectStoreUrl,
			s.manifest.Files.Clips[rand.Intn(len(s.manifest.Files.Clips))],
		)

		data := IndexRenderContext{
			InitialClipURL: clipUrl,
			InitialCaption: s.tp.Caption(""),
		}

		s.cfg.Templates.ExecuteTemplate(w, "index.html.tmpl", data)
	}
}
