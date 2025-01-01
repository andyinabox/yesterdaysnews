package domain

import "context"

type Uploader interface {
	// UploadFile() // TODO
	UploadDir(ctx context.Context, dir string) (string, error)
	PruneObjects(ctx context.Context, prefixToKeep string) ([]string, error)
	PromoteObjects(ctx context.Context, prefix string) (string, error)
}
