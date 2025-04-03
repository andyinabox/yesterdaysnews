package builder

import (
	"context"
	"log/slog"
	"path/filepath"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (b *Builder) GenerateCombinedVideo(ctx context.Context, uploadDir, modelFile string, clips []string) (string, string, error) {

	slog.Info("GenerateCombinedVideo", "uploadDir", uploadDir, "modelFile", modelFile, "clipsCount", len(clips))

	videoFile := filepath.Join(b.cfg.OutputDir, domain.VideoFileName)
	slog.Info("combining video files", "count", len(clips), "file", videoFile)
	videoFile, err := b.vp.ShuffleClipsAndCombine(ctx, clips, videoFile)
	if err != nil {
		return "", "", err
	}

	return videoFile, "", nil
}
