package objectstoreservice

import (
	"context"
	"fmt"
	"sync"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (u *Uploader) CopyObject(ctx context.Context, from, to string) (string, error) {

	err := u.osclient.CopyObject(ctx, u.cfg.ContainerName, u.cfg.ContainerName, from, to)
	if err != nil {
		return "", fmt.Errorf("error copying object %q to %q: %w", from, to, err)
	}

	return to, nil
}

func (u *Uploader) CopyObjectStream(ctx context.Context, errs chan<- domain.Error, fileKeys <-chan [2]string) <-chan string {
	stream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		close(stream)
	}

	copyObject := func(from, to string) {
		defer wg.Done()

		fileKey, err := u.CopyObject(ctx, from, to)
		if err != nil {
			errs <- domain.Err(domain.ErrTypeCopyObject, err)
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
