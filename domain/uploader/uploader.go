package uploader

import (
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/objectstoreclient"
)

// type ObjectStoreClient interface {
// 	ContainerExists(ctx context.Context, containerName string) (bool, error)
// 	CreatePublicContainer(ctx context.Context, containerName string) error
// 	UploadObject(ctx context.Context, containerName, key string, reader io.Reader, contentType string, multipart bool) error
// 	CopyObject(ctx context.Context, sourceContainer, destContainer, sourceKey, destKey string) error
// 	DeleteObject(ctx context.Context, containerName, key string) error
// 	GetObject(ctx context.Context, containerName, key string) ([]byte, error)
// 	ListObjects(ctx context.Context, containerName string) ([]string, error)
// 	ListObjectsWithPrefix(ctx context.Context, containerName string) ([]string, error)
// }

type Config struct {
	// client config
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string

	// upload config
	ContainerName string

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
