package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

func (s *Server) Captions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log.Debug("captions sse connected")

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
				// log.Debug(caption)
				fmt.Fprintf(w, "data: %s\n\n", caption)
				w.(http.Flusher).Flush()

				time.Sleep(mapCaptionToDelay(
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

func mapCaptionToDelay(cap string, minLength, maxLength int, minDelay, maxDelay float64) time.Duration {
	tokens := strings.Split(cap, " ")
	percent := float64(len(tokens)-minLength) / float64(maxLength-minLength)
	seconds := ((maxDelay - minDelay) * percent) + minDelay
	duration := time.Duration(seconds * float64(time.Second))

	// log.Debug("caption delay", "duration", duration, "percent", percent, "seconds", seconds)

	return duration
}
