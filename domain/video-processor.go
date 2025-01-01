package domain

import "context"

type BreakVideoIntoClipsResult struct {
	Files []string `json:"files"`
}

type VideoProcessor interface {
	BreakVideoIntoClips(ctx context.Context, inputPath, outDir string, minLength, maxLength int) (*BreakVideoIntoClipsResult, error)
	ShuffleClipsAndCombine(ctx context.Context, files []string, outFile string) (string, error)
}
