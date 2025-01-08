package objectstoreservice

import (
	"context"
	"errors"
)

func (s *Service) ObjectExists(ctx context.Context, fileKey string) (bool, error) {
	return false, errors.New("not implemented")
}
