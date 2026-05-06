package containerservice

import (
	"context"

	"code.andydayton.com/andy/yesterdaysnews/pkg/objectstoreclient"
)

func (s *Service) ObjectExists(ctx context.Context, fileKey string) (bool, error) {
	_, err := s.osclient.GetObject(ctx, s.cfg.ContainerName, fileKey)
	if err == nil {
		return true, nil
	}

	if err == objectstoreclient.ErrObjectDoesNotExist {
		return false, nil
	}

	return false, err
}
