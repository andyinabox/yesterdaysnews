package clipstreamer

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (s *Streamer) VideoUploadStream(ctx context.Context, errs chan<- domain.StreamErr, clipFiles <-chan string) <-chan string {
	stream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		close(stream)
	}

	uploadVideo := func(filePath string) {
		defer wg.Done()

		var err error

		fileKey := strings.Replace(filePath, s.cfg.OutputDir, s.cfg.FileUploadDir, 1)

		fileKey, err = s.up.UploadFile(ctx, filePath, fileKey, false)
		if err != nil {
			errs <- NewErr(domain.StreamErrTODO, fmt.Errorf("error uploading file %q as %q; %w", filePath, fileKey, err))
			return
		}

		stream <- fileKey
	}

	go func() {
		defer cleanup()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				filePath, open := <-clipFiles

				if !open {
					return
				}

				wg.Add(1)
				go uploadVideo(filePath)
			}
		}
	}()

	return stream
}
