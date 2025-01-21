package server

import (
	"fmt"
	"math/rand"
	"net/http"
)

type IndexRenderContext struct {
	PageTitle      string
	InitialClipURL string
}

func (s *Server) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		clipUrl := fmt.Sprintf(
			"%s/%s/%s",
			s.cfg.CDNUrl,
			s.buildID,
			s.manifest.Files.Clips[rand.Intn(len(s.manifest.Files.Clips))],
		)

		title := "yesterday's news"

		if s.manifest != nil {
			title = s.manifest.ContentDate.Format("Monday, January 2, 2006")
		}

		data := IndexRenderContext{
			PageTitle:      title,
			InitialClipURL: clipUrl,
		}

		s.cfg.Templates.ExecuteTemplate(w, "index.html.tmpl", data)
	}
}
