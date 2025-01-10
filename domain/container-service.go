package domain

import "context"

type ContainerService interface {
	GetObject(ctx context.Context, fileKey string) ([]byte, error)
	ListObjectsInDir(ctx context.Context, dirName string) ([]string, error)
	ListPrefixes(context.Context) ([]string, error)
	UploadFile(ctx context.Context, filePath, fileKey, contentType string, multipart bool) (string, error)
	UploadFileStream(ctx context.Context, errs chan<- Error, filePaths <-chan [2]string, contentType string, multipart bool) <-chan string
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
