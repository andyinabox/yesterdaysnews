package clipstreamer

import (
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

type Config struct {
	VideoDate                time.Time
	DownloadCountPerPlaylist int
	DownloadDir              string
	OutputDir                string
}

type Streamer struct {
	dl  domain.Downloader
	cfg *Config
}

func New(dl domain.Downloader, cfg *Config) *Streamer {
	return &Streamer{dl, cfg}
}
