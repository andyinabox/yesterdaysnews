package builder

import (
	"context"
	"log/slog"
	"path/filepath"
)

func (b *Builder) GenerateCombinedVideo(ctx context.Context, uploadDir, modelFile string, clips []string) (string, string, error) {

	videoFile := filepath.Join(b.cfg.OutputDir, "video.mp4")
	slog.Info("combining video files", "count", len(clips), "file", videoFile)
	videoFile, err := b.vp.ShuffleClipsAndCombine(ctx, clips, videoFile)
	if err != nil {
		return "", "", err
	}

	return videoFile, "", nil
}
