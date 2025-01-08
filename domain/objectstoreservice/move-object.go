package objectstoreservice

import (
	"context"
	"fmt"
	"sync"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (u *Uploader) MoveObject(ctx context.Context, from, to string) (string, error) {

	_, err := u.CopyObject(ctx, from, to)
	if err != nil {
		return "", fmt.Errorf("error copying object %q to %q: %w", from, to, err)
	}

	err = u.osclient.DeleteObject(ctx, u.cfg.ContainerName, from)
	if err != nil {
		return "", fmt.Errorf("error deleting object %q: %w", from, err)
	}

	return to, nil
}

func (u *Uploader) MoveObjectStream(ctx context.Context, errs chan<- domain.Error, fileKeys <-chan [2]string) <-chan string {
	stream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		close(stream)
	}

	moveObject := func(from, to string) {
		defer wg.Done()

		fileKey, err := u.MoveObject(ctx, from, to)
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
				keys, closed := <-fileKeys

				if closed {
					return
				}

				wg.Add(1)
				go moveObject(keys[0], keys[1])
			}
		}
	}()

	return stream
}
