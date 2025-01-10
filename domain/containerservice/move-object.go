package containerservice

import (
	"context"
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (s *Service) MoveObject(ctx context.Context, from, to string) (string, error) {

	_, err := s.CopyObject(ctx, from, to)
	if err != nil {
		return "", fmt.Errorf("error copying object %q to %q: %w", from, to, err)
	}

	_, err = s.DeleteObject(ctx, from)
	if err != nil {
		return "", fmt.Errorf("error deleting object %q: %w", from, err)
	}

	return to, nil
}

func (s *Service) MoveObjectStream(ctx context.Context, errs chan<- domain.Error, fileKeys <-chan [2]string) <-chan string {
	stream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		log.Debug("closing MoveObjectStream channel")
		close(stream)
	}

	moveObject := func(from, to string) {
		defer wg.Done()

		log.Debugf("move object stream %q to %q", from, to)
		fileKey, err := s.MoveObject(ctx, from, to)
		if err != nil {
			errs <- domain.Err(domain.ErrTypeMoveObject, err)
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
				keys, open := <-fileKeys

				if !open {
					log.Debug("MoveObjectStream input channel was closed")
					return
				}

				wg.Add(1)
				go moveObject(keys[0], keys[1])
			}
		}
	}()

	return stream
}
