package objectstoreservice

import (
	"context"
	"fmt"
	"os"
)

func (s *Service) UploadFile(ctx context.Context, filePath, fileKey, contentType string, multipart bool) (string, error) {

	exists, err := s.osclient.ContainerExists(ctx, s.cfg.ContainerName)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", ErrContainerDoesNotExist
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file %q: %w", filePath, err)
	}

	return s.osclient.UploadFile(
		ctx,
		s.cfg.ContainerName,
		fileKey,
		file,
		contentType,
		multipart,
	)
}
