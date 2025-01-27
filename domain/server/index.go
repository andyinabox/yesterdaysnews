package server

import (
	"fmt"
	"html/template"
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
	PageTitle      string
	InitialClipURL string
	MetaComment    template.HTML
	AboutContent   template.HTML
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

		metaComment := template.HTML(fmt.Sprintf(
			metaCommentTmpl,
			s.manifest.BuildDate,
			s.manifest.ContentDate,
			s.manifest.ID,
		))

		data := IndexRenderContext{
			PageTitle:      title,
			InitialClipURL: clipUrl,
			MetaComment:    metaComment,
			AboutContent:   template.HTML(s.cfg.AboutContent),
		}

		s.cfg.Templates.ExecuteTemplate(w, "index.html.tmpl", data)
	}
}
