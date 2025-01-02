package clipstreamer

import (
	"context"
	"sync"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (s *Streamer) VideoDownloadStream(ctx context.Context, errs chan<- domain.StreamErr, ids <-chan string) <-chan string {
	downloadPaths := make(chan string)

	var wg sync.WaitGroup

	// wait for existing
	cleanup := func() {
		log.Info("cleaning up video download stream")
		wg.Wait()
		close(downloadPaths)
	}

	// video download func
	// not checking for ctx.Done() here,
	// that should be done further down the chain
	downloadVideo := func(id string) {
		defer wg.Done()
		log.Info("Download video", "id", id)
		path, err := s.dl.DownloadVideo(ctx, id, s.cfg.DownloadDir)

		// handle error
		if err != nil {
			errs <- NewErr(domain.StreamErrDownloadVideo, err)
			return
		}

		// add path to stream
		downloadPaths <- path
	}

	// main goroutine
	go func() {
		defer cleanup()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				id, open := <-ids

				if !open {
					return
				}

				wg.Add(1)
				go downloadVideo(id)
			}
		}
	}()

	return downloadPaths
}
