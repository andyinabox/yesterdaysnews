package domain

import (
	"context"
	"time"
)

type Builder interface {
	Run(ctx context.Context, playlistIds ...string) error
	BuildVideoClips(ctx context.Context, date time.Time, uploadDir string, ids []string) []string
	BuildModel(ctx context.Context, uploadDir string, corpi []Corpus) (string, error)
	UploadManifest(ctx context.Context, uploadDir string) (string, error)
	Promote(context.Context) error
}
