package uploader

import "context"

type ObjectStoreClient interface {
	CreateBucket(ctx context.Context, bucketName string) error
	ListBuckets(ctx context.Context, prefix string) ([]string, error)
	UploadFile(ctx context.Context, bucketName, fileKey, inFilepath string) (string, error)
	UploadFileMultipart(ctx context.Context, bucketName, fileKey, inFilepath string) (string, error)
}

type Config struct {
	// upload config
	BucketNameBase string

	// manifest config
	VideoFileName    string
	SubsFileName     string
	CombinedFileName string
	ModelFileName    string
	ClipsDirName     string
}

type Uploader struct {
	osclient ObjectStoreClient
	cfg      *Config
}

func New(osclient ObjectStoreClient, cfg *Config) *Uploader {
	return &Uploader{osclient, cfg}
}
