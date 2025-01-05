package clipstreamer

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (s *Streamer) VideoUploadStream(ctx context.Context, clipFiles <-chan string) <-chan string {
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

		log.Infof("uploading video %q as %q", filePath, fileKey)

		fileKey, err = s.up.UploadFile(ctx, filePath, fileKey, "video/webm", false)
		if err != nil {
			s.error(domain.ErrTypeUploadVideo, fmt.Errorf("error uploading video %q as %q; %w", filePath, fileKey, err))
			return
		}

		log.Infof("finished uploading video %q", fileKey)

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
