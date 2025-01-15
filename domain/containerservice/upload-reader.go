package containerservice

import (
	"context"
	"io"
)

func (s *Service) UploadReader(ctx context.Context, r io.Reader, fileKey, contentType string, multipart bool) (string, error) {

	exists, err := s.osclient.ContainerExists(ctx, s.cfg.ContainerName)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", ErrContainerDoesNotExist
	}

	return s.osclient.UploadPublicFile(
		ctx,
		s.cfg.ContainerName,
		fileKey,
		r,
		contentType,
		multipart,
	)
}
