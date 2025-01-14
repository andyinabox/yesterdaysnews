package domain

import "time"

type Manifest struct {
	Date  time.Time     `json:"date"`
	ID    string        `json:"id"`
	Files ManifestFiles `json:"files"`
}

type ManifestFiles struct {
	ModelFile string   `json:"model"`
	Clips     []string `json:"clips"`
}
