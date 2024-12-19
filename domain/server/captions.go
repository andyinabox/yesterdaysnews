package server

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

func (s *Server) Captions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log.Debug("got livereload handshake from client")

		// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for {
			select {

			case <-ctx.Done():
				return
			default:
				delay := time.Duration(1.0+(3.0*rand.Float32())) * time.Second
				time.Sleep(delay)
				caption := s.cg.Caption()
				log.Debug(caption)
				fmt.Fprintf(w, "data: %s\n\n", caption)
				w.(http.Flusher).Flush()
			}
		}
	}
}
