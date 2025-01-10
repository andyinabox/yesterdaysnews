package containerservice

import (
	"context"
)

func (s *Service) GetObject(ctx context.Context, fileKey string) ([]byte, error) {
	return s.osclient.GetObject(ctx, s.cfg.ContainerName, fileKey)
}
