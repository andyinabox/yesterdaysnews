package builder

import (
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/objectstoreservice"
	"gitlab.com/andyinabox/yesterdaysnews/domain/videoprocessor"
	"gitlab.com/andyinabox/yesterdaysnews/domain/youtubeservice"
)

type Builder struct {
	yt   domain.YouTubeService
	vp   domain.VideoProcessor
	os   domain.ObjectStoreService
	eh   domain.ErrorHandler
	errs chan<- domain.Error
	cfg  *Config
}

func New(cfg *Config, eh domain.ErrorHandler) *Builder {

	// not sure if there's a better way to do this,
	// allows us to take errors as domain.Error
	// rather than errorhandler.Error
	errs := make(chan domain.Error)
	go func() {
		defer close(errs)
		for err := range errs {
			eh.Channel() <- err
		}
	}()

	yt := youtubeservice.New(&youtubeservice.Config{
		GoogleAPIKey: cfg.GoogleAPIKey,
	})

	vp := videoprocessor.New(&videoprocessor.Config{})

	os := objectstoreservice.New(&objectstoreservice.Config{
		S3Endpoint:    cfg.S3Endpoint,
		S3AccessKey:   cfg.S3AccessKey,
		S3SecretKey:   cfg.S3SecretKey,
		ContainerName: cfg.ObjectStoreContainerName,
	})

	return &Builder{
		yt:   yt,
		vp:   vp,
		os:   os,
		eh:   eh,
		errs: errs,
		cfg:  cfg,
	}
}
