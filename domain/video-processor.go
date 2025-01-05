package domain

import (
	"context"
	"time"
)

type BreakVideoIntoClipsResult struct {
	Files []string `json:"files"`
}

type VideoEdit interface {
	Start() time.Duration
	End() time.Duration
	Duration() time.Duration
}

type VideoProcessor interface {
	GetVideoEditPoints(ctx context.Context, videoFile string, minClipLength, maxClipLength time.Duration) ([]VideoEdit, error)
	CutVideo(ctx context.Context, inFile, outFile string, edit VideoEdit) (string, error)
	// BreakVideoIntoClips(ctx context.Context, inputPath, outDir string, minLength, maxLength int) (*BreakVideoIntoClipsResult, error)
	ShuffleClipsAndCombine(ctx context.Context, files []string, outFile string) (string, error)
}
