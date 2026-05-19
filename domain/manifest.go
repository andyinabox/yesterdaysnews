package domain

import "time"

type Manifest struct {
	BuildDate   time.Time     `json:"buildDate"`
	ContentDate time.Time     `json:"contentDate"`
	ID          string        `json:"id"`
	Files       ManifestFiles `json:"files"`
	Sources     []SourceStat  `json:"sources,omitempty"`
}

type ManifestFiles struct {
	PosterImageFile string   `json:"posterImage"`
	ModelFile       string   `json:"model"`
	VideoFile       string   `json:"video"`
	SubtitlesFile   string   `json:"subtitles"`
	Clips           []string `json:"clips"`
}

type SourceStat struct {
	Name            string  `json:"name"`
	PlaylistID      string  `json:"playlistId"`
	DurationSeconds float64 `json:"durationSeconds"`
	Percentage      float64 `json:"percentage"`
	ClipCount       int     `json:"clipCount"`
}
