package containerservice

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
)

func (s *Service) UploadFile(ctx context.Context, filePath, fileKey, contentType string, multipart bool) (string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file %q: %w", filePath, err)
	}

	return s.UploadReader(ctx, file, fileKey, contentType, multipart)
}

func (s *Service) UploadFileStream(ctx context.Context, errs chan<- domain.Error, filePaths <-chan [2]string, contentType string, multipart bool) <-chan string {
	stream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		log.Debug("closing upload file stream")
		close(stream)
	}

	uploadFile := func(filePath, fileKey string) {
		defer wg.Done()

		log.Infof("uploading file %q as %q", filePath, fileKey)
		fileKey, err := s.UploadFile(ctx, filePath, fileKey, contentType, multipart)
		if err != nil {
			log.Errorf("error uploading file %q as %q: %s", filePath, fileKey, err)
			errs <- errorhandler.Err(domain.ErrTypeUploadFile, err)
			return
		}
		log.Infof("finished uploading %q", fileKey)
		stream <- fileKey
	}

	go func() {
		defer cleanup()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fp, open := <-filePaths

				if !open {
					return
				}

				wg.Add(1)
				go uploadFile(fp[0], fp[1])
			}
		}
	}()

	return stream
}
