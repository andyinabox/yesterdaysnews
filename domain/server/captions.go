package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func (s *Server) Captions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		slog.Debug("captions sse connected")

		// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		var caption string

		minDelay := s.cfg.MinCaptionDelay
		maxDelay := s.cfg.MaxCaptionDelay
		minLength := s.cg.MinCaptionLength()
		maxLength := s.cg.MaxCaptionLength()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				caption = s.cg.Caption(caption)
				// slog.Debug(caption)
				fmt.Fprintf(w, "data: %s\n\n", caption)
				w.(http.Flusher).Flush()

				time.Sleep(util.MapCaptionToDelay(
					caption,
					minLength,
					maxLength,
					minDelay,
					maxDelay,
				))

			}
		}
	}
}
