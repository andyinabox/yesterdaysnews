package objectstoreservice

import (
	"context"
	"fmt"
	"sync"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (u *Uploader) DeleteObject(ctx context.Context, key string) (string, error) {
	err := u.osclient.DeleteObject(ctx, u.cfg.ContainerName, key)
	if err != nil {
		return "", fmt.Errorf("error deketing object %q: %w", key, err)
	}
	return key, nil
}

func (u *Uploader) DeleteObjectStream(ctx context.Context, errs chan<- domain.Error, fileKeys <-chan string) <-chan string {
	stream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		close(stream)
	}

	deleteObject := func(key string) {
		defer wg.Done()

		fileKey, err := u.DeleteObject(ctx, key)
		if err != nil {
			errs <- domain.Err(domain.ErrTypeDeleteObject, err)
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
				key, closed := <-fileKeys

				if closed {
					return
				}

				wg.Add(1)
				go deleteObject(key)
			}
		}
	}()

	return stream
}
