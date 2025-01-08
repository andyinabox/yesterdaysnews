package builder

import (
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/downloader"
	"gitlab.com/andyinabox/yesterdaysnews/domain/objectstoreservice"
	"gitlab.com/andyinabox/yesterdaysnews/domain/videoprocessor"
	"gitlab.com/andyinabox/yesterdaysnews/domain/youtubeservice"
)

type Builder struct {
	yt  domain.YouTubeService
	vp  domain.VideoProcessor
	up  domain.ObjectStoreService
	eh  domain.ErrorHandler
	cfg *Config
}

func New(cfg *Config, eh domain.ErrorHandler) *Builder {

	yt := youtubeservice.New(&downloader.Config{
		GoogleAPIKey: cfg.GoogleAPIKey,
	})

	vp := videoprocessor.New(&videoprocessor.Config{})

	up := objectstoreservice.New(&objectstoreservice.Config{
		S3Endpoint:  cfg.S3Endpoint,
		S3AccessKey: cfg.S3AccessKey,
		S3SecretKey: cfg.S3SecretKey,

		ContainerName: cfg.ObjectStoreContainerName,
		PrimaryDir:    cfg.ObjectStorePrimaryDir,
	})

	return &Builder{
		yt:  yt,
		vp:  vp,
		up:  up,
		eh:  eh,
		cfg: cfg,
	}
}
