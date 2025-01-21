package domain

import "time"

type Manifest struct {
	BuildDate   time.Time     `json:"buildDate"`
	ContentDate time.Time     `json:"contentDate"`
	ID          string        `json:"id"`
	Files       ManifestFiles `json:"files"`
}

type ManifestFiles struct {
	ModelFile string   `json:"model"`
	Clips     []string `json:"clips"`
}
