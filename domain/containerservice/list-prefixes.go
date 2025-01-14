package containerservice

import "context"

func (s *Service) ListPrefixes(ctx context.Context) ([]string, error) {
	return s.osclient.ListPrefixes(ctx, s.cfg.ContainerName)
}
