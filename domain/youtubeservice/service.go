package youtubeservice

import (
	"time"

	"code.andydayton.com/andy/yesterdaysnews/pkg/youtubeapi"
	"code.andydayton.com/andy/yesterdaysnews/pkg/youtubedownloader"
)

const VideoFormatString = "bv*[width<=1280][aspect_ratio>1]"
const VideoSubFormat = "vtt"

type Config struct {
	GoogleAPIKey        string
	BinPathYTDLP        string
	MaxPlaylistRequests int
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
