package clipstreamer

import (
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/errorhandler"
)

type Config struct {
	VideoDate                time.Time
	DownloadCountPerPlaylist int
	DownloadDir              string
	OutputDir                string
	ClipsDir                 string
	MinClipLength            time.Duration
	MaxClipLength            time.Duration
	FileUploadDir            string
}

type Streamer struct {
	dl   domain.Downloader
	vp   domain.VideoProcessor
	up   domain.Uploader
	errs chan<- errorhandler.Error
	cfg  *Config
}

func New(dl domain.Downloader, vp domain.VideoProcessor, up domain.Uploader, errs chan<- errorhandler.Error, cfg *Config) *Streamer {
	return &Streamer{dl, vp, up, errs, cfg}
}
