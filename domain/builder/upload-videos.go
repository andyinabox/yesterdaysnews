package builder

import (
	"context"
	"log/slog"
	"path/filepath"
	"strings"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
)

func mimeTypeFromExt(path string) string {
	switch filepath.Ext(path) {
	case ".mp4":
		return "video/mp4"
	default:
		return "video/webm"
	}
}

func (b *Builder) UploadVideos(ctx context.Context, uploadDir string, clips <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for clipPath := range clips {
			select {
			case <-ctx.Done():
				return
			default:
			}

			fileKey := strings.Replace(clipPath, b.cfg.OutputDir, uploadDir, 1)
			contentType := mimeTypeFromExt(clipPath)

			key, err := b.cs.UploadFile(ctx, clipPath, fileKey, contentType, false)
			if err != nil {
				b.errs <- errorhandler.Err(domain.ErrTypeUploadFile, err)
				continue
			}

			slog.Debug("uploaded clip", "key", key)

			// remove upload dir to get relative paths from upload root
			relPath := strings.Replace(key, uploadDir+"/", "", 1)
			out <- relPath
		}
	}()

	return out
}
