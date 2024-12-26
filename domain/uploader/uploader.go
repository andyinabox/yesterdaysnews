package uploader

import (
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/objectstoreclient"
)

type Config struct {
	// client config
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string

	// upload config
	BucketNameBase string

	// manifest config
	VideoFileName string
	// SubsFileName     string
	// CombinedFileName string
	ModelFileName string
	ClipsDirName  string
}

type Uploader struct {
	osclient *objectstoreclient.Client
	cfg      *Config
}

func New(cfg *Config) *Uploader {
	return &Uploader{
		osclient: objectstoreclient.New(&objectstoreclient.Config{
			Endpoint:  cfg.S3Endpoint,
			AccessKey: cfg.S3AccessKey,
			SecretKey: cfg.S3SecretKey,
		}),
		cfg: cfg,
	}
}
