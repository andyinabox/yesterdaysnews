package server

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
)

type clipsResponse struct {
	Clips []string `json:"clips"`
}

func (s *Server) Clips() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clips := make([]string, len(s.manifest.Files.Clips))

		for i, c := range s.manifest.Files.Clips {
			clips[i] = fmt.Sprintf("%s/%s", s.cfg.ObjectStoreUrl, c)
		}

		rand.Shuffle(len(clips), func(i, j int) {
			clips[i], clips[j] = clips[j], clips[i]
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clipsResponse{
			Clips: clips,
		})
	}
}
