package containerservice

import (
	"context"
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
)

func (s *Service) DeleteObject(ctx context.Context, key string) (string, error) {
	err := s.osclient.DeleteObject(ctx, s.cfg.ContainerName, key)
	if err != nil {
		return "", fmt.Errorf("error deketing object %q: %w", key, err)
	}
	return key, nil
}

func (s *Service) DeleteObjectStream(ctx context.Context, errs chan<- domain.Error, fileKeys <-chan string) <-chan string {
	stream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		close(stream)
	}

	deleteObject := func(key string) {
		defer wg.Done()

		log.Debugf("delete %q from stream", key)
		fileKey, err := s.DeleteObject(ctx, key)
		if err != nil {
			errs <- errorhandler.Err(domain.ErrTypeDeleteObject, err)
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
				key, open := <-fileKeys

				if !open {
					return
				}

				wg.Add(1)
				go deleteObject(key)
			}
		}
	}()

	return stream
}
