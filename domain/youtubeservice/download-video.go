package youtubeservice

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/youtubedownloader"
)

func (s *Service) DownloadVideo(ctx context.Context, id, outDir string) (string, error) {
	output := filepath.Join(outDir, "%(id)s.%(ext)s")

	video, err := s.ytdl.DownloadVideo(ctx, id, youtubedownloader.Request{
		Format:        VideoFormatString,
		WriteAutoSubs: true,
		SubFormat:     VideoSubFormat,
		Output:        output,
	})
	if err != nil {
		return "", err
	}

	return video.Filename, nil
}

func (s *Service) DownloadVideoStream(ctx context.Context, errs chan<- domain.Error, ids <-chan string, outDir string) <-chan string {
	downloadPaths := make(chan string)

	var wg sync.WaitGroup

	// wait for existing
	cleanup := func() {
		log.Debug("cleaning up video download stream")
		wg.Wait()
		close(downloadPaths)
	}

	// video download func
	// that should be done further down the chain
	downloadVideo := func(id string) {
		defer wg.Done()

		log.Info("Download video", "id", id)
		path, err := s.DownloadVideo(ctx, id, outDir)

		// handle error
		if err != nil {
			errs <- errorhandler.Err(domain.ErrTypeDownloadVideo, err)
			return
		}

		// add path to stream
		downloadPaths <- path
	}

	// main goroutine
	go func() {
		defer cleanup()
		for id := range ids {
			select {
			case <-ctx.Done():
				return
			default:
				wg.Add(1)
				go downloadVideo(id)
			}
		}
	}()

	return downloadPaths
}
