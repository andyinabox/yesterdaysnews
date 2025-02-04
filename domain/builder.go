package domain

import (
	"context"
	"time"
)

type BuildPhase string

const (
	BuildPhaseAll            BuildPhase = "all"
	BuildPhaseSetup          BuildPhase = "setup"
	BuildPhaseDownloadVideos BuildPhase = "download-videos"
	BuildPhaseCutVideos      BuildPhase = "cut-videos"
	BuildPhaseUploadVideos   BuildPhase = "upload-videos"
	BuildPhaseExtractImages  BuildPhase = "extract-images"
	BuildPhasePosterImage    BuildPhase = "poster-image"
	BuildPhaseModel          BuildPhase = "model"
	BuildPhaseManifest       BuildPhase = "manifest"
	BuildPhasePromote        BuildPhase = "promote"
	BuildPhaseCleanup        BuildPhase = "cleanup"
)

type Builder interface {
	// Run will execute the below subcommands in sequence
	Run(ctx context.Context) error

	// these are mainly exported so they can be tested with subcommands
	Setup(ctx context.Context) error
	DownloadVideos(ctx context.Context, date time.Time) <-chan string
	CutVideos(ctx context.Context, paths <-chan string) <-chan string
	UploadVideos(ctx context.Context, uploadDir string, clips <-chan string) <-chan string
	ExtractImages(ctx context.Context, clips <-chan string) <-chan string
	// VideoClips(ctx context.Context, date time.Time, uploadDir string) ([]string, error)
	PosterImage(ctx context.Context, uploadDir string, clips <-chan string) (string, error)
	Model(ctx context.Context, uploadDir string, corpi []Corpus) (string, error)
	Manifest(ctx context.Context, uploadDir string, manifest *Manifest) (string, error)
	Promote(ctx context.Context, uploadDir string) error
	Cleanup(ctx context.Context, toKeep int) ([]string, error)
}
