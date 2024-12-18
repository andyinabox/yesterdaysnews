package downloader

import (
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/youtubeapi"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/youtubedownloader"
)

const VideoFormatString = "bv[ext=mp4][aspect_ratio>1]"

type Config struct {
	GoogleAPIKey string
	BinPathYTDLP string
}

type Downloader struct {
	ytapi *youtubeapi.Client
	ytdl  *youtubedownloader.Client
	cfg   *Config
}

func New(cfg *Config) *Downloader {
	return &Downloader{
		ytapi: youtubeapi.New(cfg.GoogleAPIKey),
		ytdl:  youtubedownloader.New(cfg.BinPathYTDLP),
	}
}
