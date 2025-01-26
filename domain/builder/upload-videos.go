package builder

import (
	"context"
	"strings"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/streams"
)

func (b *Builder) UploadVideos(ctx context.Context, uploadDir string, clips <-chan string) <-chan string {
	// take clip paths amd make stream of upload paths
	uploadPaths := streams.StringStreamTo2StringSliceStream(ctx, clips, func(s string) [2]string {
		// first string is the file name, second is the the object key
		return [2]string{s, strings.Replace(s, b.cfg.OutputDir, uploadDir, 1)}
	})

	// upload files
	clipUploads := b.cs.UploadFileStream(ctx, b.errs, uploadPaths, "video/webm", false)

	// remove upload dir to get relative paths from upload root
	clipPathsStream := streams.StringTransformStream(ctx, clipUploads, func(s string) string {
		return strings.Replace(s, uploadDir+"/", "", 1)
	})

	return clipPathsStream
}
