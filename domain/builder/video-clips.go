package builder

// import (
// 	"context"
// 	"path/filepath"
// 	"strings"
// 	"sync"
// 	"time"

// 	"github.com/charmbracelet/log"
// 	"gitlab.com/andyinabox/yesterdaysnews/domain"
// 	"gitlab.com/andyinabox/yesterdaysnews/pkg/streams"
// 	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
// )

// func (b *Builder) VideoClips(ctx context.Context, date time.Time, uploadDir string) ([]string, error) {

// 	// get stream of ids
// 	ids := b.videoIDStream(ctx, date, b.cfg.PlaylistIDs)

// 	// get stream of downloaded video paths
// 	paths := b.yt.DownloadVideoStream(ctx, b.errs, ids, b.cfg.OutputDir)

// 	// get stream of video clip paths
// 	clips := b.videoCutStream(ctx, paths)

// 	if b.cfg.SkipUpload {
// 		log.Info("SkipUpload is true, skipping video uploads")

// 		clipPaths := streams.StringTransformStream(ctx, clips, func(s string) string {
// 			return strings.TrimPrefix(s, b.cfg.OutputDir+"/")
// 		})

// 		return streams.StringSlice(ctx, clipPaths), nil
// 	}

// 	// take clip paths amd make stream of upload paths
// 	uploadPaths := streams.StringStreamTo2StringSliceStream(ctx, clips, func(s string) [2]string {
// 		// first string is the file name, second is the the object key
// 		return [2]string{s, strings.Replace(s, b.cfg.OutputDir, uploadDir, 1)}
// 	})

// 	// upload files
// 	clipUploads := b.cs.UploadFileStream(ctx, b.errs, uploadPaths, "video/webm", false)

// 	// remove upload dir to get relative paths from upload root
// 	clipPathsStream := streams.StringTransformStream(ctx, clipUploads, func(s string) string {
// 		return strings.Replace(s, uploadDir+"/", "", 1)
// 	})

// 	// wait for all clips to be complete and return a string slice
// 	return streams.StringSlice(ctx, clipPathsStream), nil
// }
