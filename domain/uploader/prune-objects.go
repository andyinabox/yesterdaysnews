package uploader

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/objectstoreclient"
)

func (u *Uploader) PruneObjects(ctx context.Context, containerToKeep string) ([]string, error) {

	// containers, err := u.osclient.ListBuckets(ctx, &objectstoreclient.ListBucketsRequest{
	// 	Prefix: fmt.Sprintf("%s-", u.cfg.BucketName),
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// var removed []string

	// for _, containerName := range containers {
	// 	if containerName != containerToKeep {
	// 		log.Debugf("deleting container %q", containerName)
	// 		err = u.osclient.DeleteBucket(ctx, containerName, true)
	// 		if err != nil {
	// 			return removed, fmt.Errorf("error deleting container %q: %w", containerName, err)
	// 		}
	// 		removed = append(removed, containerName)
	// 	}
	// }

	// return removed, nil
	return nil, errors.New("not implemented")
}

func (u *Uploader) deleteObjectsWithPrefix(ctx context.Context, prefix string) ([]string, error) {
	objectKeys, err := u.osclient.ListObjects(ctx, u.cfg.ContainerName, &objectstoreclient.ListObjectsRequest{
		Prefix: prefix,
	})
	if err != nil {
		return nil, fmt.Errorf("error listing objects: %w", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	deleted := []string{}

	for _, key := range objectKeys {
		wg.Add(1)
		go func() {
			defer wg.Done()

			log.Debugf("deleting %q", key)
			err := u.osclient.DeleteObject(ctx, u.cfg.ContainerName, key)
			if err != nil {
				log.Errorf("error deleting %q: %s", key, err)
			}
			log.Debugf("finished deleting %q", key)

			mu.Lock()
			deleted = append(deleted, key)
			mu.Unlock()
		}()
	}

	wg.Wait()

	return deleted, nil
}
