package containerservice

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"code.andydayton.com/andy/yesterdaysnews/domain"
	"code.andydayton.com/andy/yesterdaysnews/domain/errorhandler"
	"code.andydayton.com/andy/yesterdaysnews/pkg/objectstoreclient"
)

func (s *Service) ListObjectsWithPrefix(ctx context.Context, prefix string) ([]string, error) {

	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	req := &objectstoreclient.ListObjectsRequest{
		Prefix: prefix,
	}

	result, err := s.osclient.ListObjects(ctx, s.cfg.ContainerName, req)
	if err != nil {
		return nil, fmt.Errorf("error listing object with prefix %q: %w", req.Prefix, err)
	}

	return result, nil
}

func (s *Service) ListObjectsWithPrefixStream(ctx context.Context, errs chan<- domain.Error, prefixes <-chan string) <-chan string {
	stream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		close(stream)
	}

	getObjectsWithPrefix := func(prefix string) {
		defer wg.Done()

		slog.Debug("get objects with prefix", "prefix", prefix)

		objects, err := s.ListObjectsWithPrefix(ctx, prefix)
		if err != nil {
			errs <- errorhandler.Err(domain.ErrTypeDeleteObject, err)
			return
		}

		slog.Debug("found objects with prefix", "count", len(objects), "prefix", prefix)

		for _, object := range objects {
			slog.Debug("object with prefix", "object", object, "prefix", prefix)
			stream <- object
		}
	}

	go func() {
		defer cleanup()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				prefix, open := <-prefixes

				if !open {
					return
				}

				wg.Add(1)
				go getObjectsWithPrefix(prefix)

			}
		}
	}()

	return stream
}
