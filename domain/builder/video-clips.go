package builder

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func (b *Builder) VideoClips(ctx context.Context, date time.Time, uploadDir string) ([]string, error) {

	var mu sync.Mutex

	ids := b.videoIDStream(ctx, date, b.cfg.PlaylistIDs)
	paths := b.yt.DownloadVideoStream(ctx, b.errs, ids, b.cfg.OutputDir)
	clips := b.videoCutStream(ctx, paths)
	clipUploads := b.videoUploadStream(ctx, clips, uploadDir)

	log.Debug("waiting to finish converting upload filenames")
	clipPaths := []string{}

	for fileKey := range clipUploads {
		log.Debugf("converting filename for %q", fileKey)
		mu.Lock()
		clipPaths = append(clipPaths, strings.Replace(fileKey, uploadDir+"/", "", 1))
		mu.Unlock()
	}

	log.Info("done building video clips")

	// TODO: is there any way of returning errors here?
	return clipPaths, nil
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
			b.error(domain.ErrTypeCutVideo, err)
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

func (b *Builder) videoUploadStream(ctx context.Context, filePaths <-chan string, uploadDir string) <-chan string {

	// create upload fileKeys and send on 2-string channel
	uploadPaths := make(chan [2]string)
	go func() {
		defer close(uploadPaths)
		for filePath := range filePaths {
			fileKey := strings.Replace(filePath, b.cfg.OutputDir, uploadDir, 1)
			log.Debugf("add upload path for %q, %q", filePath, fileKey)
			uploadPaths <- [2]string{filePath, fileKey}
		}
		log.Debug("finished translating upload paths, closing")
	}()

	return b.cs.UploadFileStream(ctx, b.errs, uploadPaths, "video/webm", false)
}
