package domain

import (
	"context"
	"time"
)

type BuildPhase string

const (
	BuildPhaseAll        BuildPhase = "all"
	BuildPhaseSetup      BuildPhase = "setup"
	BuildPhaseVideoClips BuildPhase = "video-clips"
	BuildPhaseModel      BuildPhase = "model"
	BuildPhaseManifest   BuildPhase = "manifest"
	BuildPhasePromote    BuildPhase = "promote"
	BuildPhaseCleanup    BuildPhase = "cleanup"
)

type Builder interface {
	// Run will execute the below subcommands in sequence
	Run(ctx context.Context) error

	// these are mainly exported so they can be tested with subcommands
	Setup(ctx context.Context) error
	VideoClips(ctx context.Context, date time.Time, uploadDir string) ([]string, error)
	Model(ctx context.Context, uploadDir string, corpi []Corpus) (string, error)
	Manifest(ctx context.Context, uploadDir string, manifest *Manifest) (string, error)
	Promote(ctx context.Context, uploadDir string) error
	Cleanup(ctx context.Context, toKeep int) ([]string, error)
}
