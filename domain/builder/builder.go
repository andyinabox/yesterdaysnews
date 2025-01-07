package builder

import (
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/downloader"
	"gitlab.com/andyinabox/yesterdaysnews/domain/uploader"
	"gitlab.com/andyinabox/yesterdaysnews/domain/videoprocessor"
)

type Builder struct {
	dl  domain.Downloader
	vp  domain.VideoProcessor
	up  domain.Uploader
	eh  domain.ErrorHandler
	cfg *Config
}

func New(cfg *Config, eh domain.ErrorHandler) *Builder {

	dl := downloader.New(&downloader.Config{
		GoogleAPIKey: cfg.GoogleAPIKey,
	})

	vp := videoprocessor.New(&videoprocessor.Config{})

	up := uploader.New(&uploader.Config{
		S3Endpoint:  cfg.S3Endpoint,
		S3AccessKey: cfg.S3AccessKey,
		S3SecretKey: cfg.S3SecretKey,

		ContainerName: cfg.ObjectStoreContainerName,
		PrimaryDir:    cfg.ObjectStorePrimaryDir,
	})

	return &Builder{
		dl:  dl,
		vp:  vp,
		up:  up,
		eh:  eh,
		cfg: cfg,
	}
}
