package domain

import "time"

const (
	ManifestVideoFileName    = "yesterdays-news.mp4"
	ManifestSubsVileName     = "yesterdays-news.en.vtt"
	ManifestCombinedFileName = "yesterdays-news-cc.mp4"
	ManifestModelFileName    = "yesterdays-news.model.json"
	ManifestClipsDirName     = "clips"
)

type Manifest struct {
	Date  time.Time     `json:"date"`
	ID    string        `json:"id"`
	Files ManifestFiles `json:"files"`
}

type ManifestFiles struct {
	VideoFile string `json:"video"`
	// SubsFile     string   `json:"subs"`
	// CombinedFile string   `json:"combined"`
	ModelFile string   `json:"model"`
	Clips     []string `json:"clips"`
}
