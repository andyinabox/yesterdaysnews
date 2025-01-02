package uploader

import (
	"context"
	"errors"
)

func (u *Uploader) UploadFile(ctx context.Context, filePath, fileKey string, multipart bool) (string, error) {
	return "", errors.New("not implemented")
}
