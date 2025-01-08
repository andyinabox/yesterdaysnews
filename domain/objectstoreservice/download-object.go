package objectstoreservice

import (
	"context"
	"errors"
)

func (s *Service) DownloadObject(ctx context.Context, fileKey string) ([]byte, error) {
	return nil, errors.New("not implemented")
}
