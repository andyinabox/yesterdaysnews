package domain

import (
	"context"
	"time"
)

type Builder interface {
	Run(ctx context.Context) error

	// these are mainly exported so they can be tested with subcommands
	BuildVideoClips(ctx context.Context, date time.Time, uploadDir string) []string
	BuildModel(ctx context.Context, uploadDir string, corpi []Corpus) (string, error)
	UploadManifest(ctx context.Context, uploadDir string, manifest *Manifest) error
	Promote(ctx context.Context, uploadDir string) error
}
