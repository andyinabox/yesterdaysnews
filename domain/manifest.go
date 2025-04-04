package domain

import "time"

type Manifest struct {
	BuildDate   time.Time     `json:"buildDate"`
	ContentDate time.Time     `json:"contentDate"`
	ID          string        `json:"id"`
	Files       ManifestFiles `json:"files"`
}

type ManifestFiles struct {
	PosterImageFile string   `json:"posterImage"`
	ModelFile       string   `json:"model"`
	VideoFile       string   `json:"video"`
	SubtitlesFile   string   `json:"subtitles"`
	Clips           []string `json:"clips"`
}
