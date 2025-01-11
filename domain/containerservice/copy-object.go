package containerservice

import (
	"context"
	"fmt"
	"sync"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
)

func (s *Service) CopyObject(ctx context.Context, from, to string) (string, error) {

	err := s.osclient.CopyObject(ctx, s.cfg.ContainerName, s.cfg.ContainerName, from, to)
	if err != nil {
		return "", fmt.Errorf("error copying object %q to %q: %w", from, to, err)
	}

	return to, nil
}

func (s *Service) CopyObjectStream(ctx context.Context, errs chan<- domain.Error, fileKeys <-chan [2]string) <-chan string {
	stream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		close(stream)
	}

	copyObject := func(from, to string) {
		defer wg.Done()

		fileKey, err := s.CopyObject(ctx, from, to)
		if err != nil {
			errs <- errorhandler.Err(domain.ErrTypeCopyObject, err)
			return
		}

		stream <- fileKey

	}

	go func() {
		defer cleanup()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				keys, closed := <-fileKeys

				if closed {
					return
				}

				wg.Add(1)
				go copyObject(keys[0], keys[1])
			}
		}
	}()

	return stream
}
