package clipstreamer

import (
	"context"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (s *Streamer) VideoDownloadStream(ctx context.Context, errs chan<- domain.StreamErr, ids <-chan string) <-chan string {
	downloadPaths := make(chan string)

	// video download func
	// not checking for ctx.Done() here,
	// that should be done further down the chain
	downloadVideo := func(id string) {
		path, err := s.dl.DownloadVideo(ctx, id, s.cfg.DownloadDir)

		// handle error
		if err != nil {
			errs <- NewStreamErr(domain.StreamErrTODO, err)
			return
		}

		// add path to stream
		downloadPaths <- path
	}

	// main goroutine
	go func() {
		defer close(downloadPaths)
		for {
			select {
			case <-ctx.Done():
				return
			case id := <-ids:
				go downloadVideo(id)
			}
		}
	}()

	return downloadPaths
}
