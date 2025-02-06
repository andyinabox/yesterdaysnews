package youtubeservice

import (
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/youtubeapi"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/youtubedownloader"
)

const VideoFormatString = "bv[ext=webm][width<=1280][aspect_ratio>1]"
const VideoSubFormat = "vtt"

type Config struct {
	GoogleAPIKey        string
	BinPathYTDLP        string
	ThrottleDownloadsBy time.Duration
}

type Service struct {
	ytapi *youtubeapi.Client
	ytdl  *youtubedownloader.Client
	cfg   *Config
}

func New(cfg *Config) *Service {
	return &Service{
		ytapi: youtubeapi.New(cfg.GoogleAPIKey),
		ytdl:  youtubedownloader.New(cfg.BinPathYTDLP),
		cfg:   cfg,
	}
}
