package domain

import (
	"context"
	"time"
)

type Builder interface {
	// Run will execute the below subcommands in sequence
	Run(ctx context.Context) error

	// these are mainly exported so they can be tested with subcommands
	Setup(ctx context.Context) error
	VideoClips(ctx context.Context, date time.Time, uploadDir string) ([]string, error)
	Model(ctx context.Context, uploadDir string, corpi []Corpus) (string, error)
	Manifest(ctx context.Context, uploadDir string, manifest *Manifest) (string, error)
	Promote(ctx context.Context, uploadDir string) (string, error)
	Cleanup(ctx context.Context, toKeep string) ([]string, error)
}
