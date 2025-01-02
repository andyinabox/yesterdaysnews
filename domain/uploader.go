package domain

import "context"

// TODO: rename to ObjectStoreService
type Uploader interface {
	UploadFile(ctx context.Context, filePath, fileKey string, multipart bool) (string, error)
	UploadDir(ctx context.Context, dir string) (string, error)
	PruneObjects(ctx context.Context, prefixToKeep string) ([]string, error)
	PromoteObjects(ctx context.Context, prefix string) (string, error)
}
