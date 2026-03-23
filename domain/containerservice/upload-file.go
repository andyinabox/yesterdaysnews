package containerservice

import (
	"context"
	"fmt"
	"log/slog"
	"mime"
	"os"
	"path/filepath"
	"sync"

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

func (s *Service) UploadFileStream(ctx context.Context, errs chan<- domain.Error, filePaths <-chan [2]string, multipart bool) <-chan string {
	stream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		slog.Debug("closing upload file stream")
		close(stream)
	}

	uploadFile := func(filePath, fileKey string) {
		defer wg.Done()

		contentType := mime.TypeByExtension(filepath.Ext(filePath))
		slog.Info("uploading file", "path", filePath, "key", fileKey)
		fileKey, err := s.UploadFile(ctx, filePath, fileKey, contentType, multipart)
		if err != nil {
			slog.Error("error uploading file", "path", filePath, "key", fileKey, "error", err)
			errs <- errorhandler.Err(domain.ErrTypeUploadFile, err)
			return
		}
		slog.Info("finished uploading", "key", fileKey)
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
