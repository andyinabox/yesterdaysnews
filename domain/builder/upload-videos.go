package builder

import (
	"context"
	"strings"

	"code.andydayton.com/andy/yesterdaysnews/pkg/streams"
)

func (b *Builder) UploadVideos(ctx context.Context, uploadDir string, clips <-chan string) <-chan string {
	filePathsStream := streams.StringStreamTo2StringSliceStream(ctx, clips, func(clipPath string) [2]string {
		fileKey := strings.Replace(clipPath, b.cfg.OutputDir, uploadDir, 1)
		return [2]string{clipPath, fileKey}
	})

	uploadedStream := b.cs.UploadFileStream(ctx, b.errs, filePathsStream, false)

	return streams.StringTransformStream(ctx, uploadedStream, func(key string) string {
		return strings.Replace(key, uploadDir+"/", "", 1)
	})
}
