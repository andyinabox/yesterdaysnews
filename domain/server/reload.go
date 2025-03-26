package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (s *Server) Reload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		slog.Debug("reload sse connected")

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
			case <-s.reload:
				slog.Debug("request browser reload")
				fmt.Fprintf(w, "data: reload\n\n")
				w.(http.Flusher).Flush()
			}
		}
	}
}
