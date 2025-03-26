package server

import (
	"fmt"
	"html/template"
	"log/slog"
	"math/rand"
	"net/http"
)

const metaCommentTmpl = `
<!--
  BuildDate: %s
  ContentDate: %s
  BuildID: %s
-->
`

type IndexRenderContext struct {
	RenderContext
	InitialClipURL          string
	MetaComment             template.HTML
	FetchClipsWhenLowerThan int
}

func (s *Server) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		clipUrl := fmt.Sprintf(
			"%s/%s/%s",
			s.cfg.CDNUrl,
			s.buildID,
			s.manifest.Files.Clips[rand.Intn(len(s.manifest.Files.Clips))],
		)

		metaComment := template.HTML(fmt.Sprintf(
			metaCommentTmpl,
			s.manifest.BuildDate,
			s.manifest.ContentDate,
			s.manifest.ID,
		))

		data := IndexRenderContext{
			RenderContext:           s.renderContext(),
			InitialClipURL:          clipUrl,
			MetaComment:             metaComment,
			FetchClipsWhenLowerThan: s.cfg.FetchClipsWhenLowerThan,
		}

		if s.manifest != nil {
			data.PageTitle = s.manifest.ContentDate.Format("Monday, January 2, 2006")
		}

		err := s.cfg.Templates.ExecuteTemplate(w, "index.html.tmpl", data)
		if err != nil {
			slog.Error("error rendering index page", "error", err)
		}
	}
}
