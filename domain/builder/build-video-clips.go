package builder

import (
	"context"
	"strings"
	"time"
)

func (b *Builder) BuildVideoClips(ctx context.Context, date time.Time, uploadDir string, playlistIDs []string) []string {

	ids := b.videoIDStream(ctx, date, playlistIDs)
	paths := b.videoDownloadStream(ctx, ids)
	clips := b.videoCutStream(ctx, paths)
	clipUploads := b.videoUploadStream(ctx, uploadDir, clips)

	clipPaths := []string{}
	for fileKey := range clipUploads {
		clipPaths = append(clipPaths, strings.Replace(fileKey, uploadDir+"/", "", 1))
	}

	return clipPaths
}
