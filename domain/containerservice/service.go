package containerservice

import (
	"gitlab.com/andyinabox/yesterdaysnews/pkg/objectstoreclient"
)

type Config struct {
	// client config
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string

	// upload config
	ContainerName string
}

type Service struct {
	osclient *objectstoreclient.Client
	cfg      *Config
}

func New(cfg *Config) *Service {
	return &Service{
		osclient: objectstoreclient.New(&objectstoreclient.Config{
			Endpoint:  cfg.S3Endpoint,
			AccessKey: cfg.S3AccessKey,
			SecretKey: cfg.S3SecretKey,
		}),
		cfg: cfg,
	}
}
