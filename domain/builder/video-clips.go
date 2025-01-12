package builder

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/streams"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func (b *Builder) VideoClips(ctx context.Context, date time.Time, uploadDir string) ([]string, error) {

	// get stream of ids
	ids := b.videoIDStream(ctx, date, b.cfg.PlaylistIDs)

	// get stream of downloaded video paths
	paths := b.yt.DownloadVideoStream(ctx, b.errs, ids, b.cfg.OutputDir)

	// get stream of video clip paths
	clips := b.videoCutStream(ctx, paths)

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

	// wait for all clips to be complete and return a string slice
	return streams.StringSlice(ctx, clipPathsStream), nil
}

func (b *Builder) videoCutStream(ctx context.Context, videoFiles <-chan string) <-chan string {
	clipStream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		log.Debug("closing video clip stream")
		close(clipStream)
	}

	cutVideo := func(filePath string) {
		defer wg.Done()

		editPoints, err := b.vp.GetVideoEditPoints(ctx, filePath, util.Seconds(b.cfg.MinClipLengthSeconds), util.Seconds(b.cfg.MaxClipLengthSeconds))

		if err != nil {
			b.eh.Add(domain.ErrTypeCutVideo, err)
			return
		}

		log.Infof("cutting %q into %d clips", filePath, len(editPoints))
		for clip := range b.vp.CutVideoStream(ctx, b.errs, filePath, editPoints) {
			log.Infof("finished cutting %q", clip)
			clipStream <- clip
		}
		log.Info("exiting cutVideo loop")
	}

	go func() {
		defer cleanup()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				filePath, open := <-videoFiles

				// stop and cleanup if the channel is closed
				if !open {
					return
				}

				wg.Add(1)
				log.Infof("start cutting %q", filePath)
				go cutVideo(filePath)
			}
		}
	}()

	return clipStream
}

func (b *Builder) videoIDStream(ctx context.Context, date time.Time, playlistIDs []string) <-chan string {

	stream := make(chan string)

	var wg sync.WaitGroup

	// combine multiple playlist streams into single stream
	wg.Add(len(playlistIDs))
	for _, playlistID := range playlistIDs {
		go func() {
			defer wg.Done()
			log.Infof("getting video id stream for %q", playlistID)
			playlistIDs := b.yt.GetPlaylistVideoIDStream(ctx, b.errs, playlistID, date, b.cfg.DownloadCountPerPlaylist)
			for id := range playlistIDs {
				log.Infof("got new video ID: %s", id)
				stream <- id
			}
		}()
	}

	// close stream once all playlistIDs are gathered
	go func() {
		wg.Wait()
		log.Debug("closing video id stream")
		close(stream)
	}()

	return stream
}
