package domain

import "context"

type ObjectStoreService interface {
	UploadFile(ctx context.Context, filePath, fileKey, contentType string, multipart bool) (string, error)
	ListObjectsInDir(ctx context.Context, dirName string) ([]string, error)
	CopyObject(ctx context.Context, from, to string) (string, error)
	CopyObjectStream(context.Context, chan<- Error, <-chan [2]string) <-chan string
	MoveObject(ctx context.Context, from, to string) (string, error)
	MoveObjectStream(context.Context, chan<- Error, <-chan [2]string) <-chan string
	DeleteObject(ctx context.Context, key string) (string, error)
	DeleteObjectStream(context.Context, chan<- Error, <-chan string) <-chan string
	// UploadDir(ctx context.Context, dir string) (string, error)
	// PruneObjects(ctx context.Context, prefixToKeep string) ([]string, error)
	// PromoteObjects(ctx context.Context, prefix string) (string, error)
}
