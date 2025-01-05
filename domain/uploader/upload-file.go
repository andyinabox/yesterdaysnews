package uploader

import (
	"context"
	"fmt"
	"os"
)

func (u *Uploader) UploadFile(ctx context.Context, filePath, fileKey, contentType string, multipart bool) (string, error) {

	exists, err := u.osclient.ContainerExists(ctx, u.cfg.ContainerName)
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

	return u.osclient.UploadFile(
		ctx,
		u.cfg.ContainerName,
		fileKey,
		file,
		contentType,
		multipart,
	)
}
