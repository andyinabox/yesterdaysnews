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
	MinClipLength            time.Duration
	MaxClipLength            time.Duration
	FileUploadDir            string
}

type Streamer struct {
	dl  domain.Downloader
	vp  domain.VideoProcessor
	up  domain.Uploader
	cfg *Config
}

func New(dl domain.Downloader, vp domain.VideoProcessor, up domain.Uploader, cfg *Config) *Streamer {
	return &Streamer{dl, vp, up, cfg}
}
