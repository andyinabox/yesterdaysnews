package containerservice

import (
	"context"
	"fmt"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/objectstoreclient"
)

func (s *Service) ListObjectsInDir(ctx context.Context, dirName string) ([]string, error) {

	req := &objectstoreclient.ListObjectsRequest{
		Prefix: dirName + "/",
	}

	result, err := s.osclient.ListObjects(ctx, s.cfg.ContainerName, req)
	if err != nil {
		return nil, fmt.Errorf("error listing object with prefix %q: %w", req.Prefix, err)
	}

	return result, nil
}
