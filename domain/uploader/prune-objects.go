package uploader

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/objectstoreclient"
)

func (u *Uploader) PruneObjects(ctx context.Context, prefixToKeep string) ([]string, error) {

	prefixes, err := u.osclient.ListPrefixes(ctx, u.cfg.ContainerName)
	if err != nil {
		return nil, err
	}

	allDeleted := []string{}

	for _, prefix := range prefixes {
		// skip primary dir
		if strings.HasPrefix(prefix, u.cfg.PrimaryDir) {
			continue
		}
		// skip keep prefix
		if prefixToKeep != "" && strings.HasPrefix(prefix, prefixToKeep) {
			continue
		}

		deleted, err := u.deleteObjectsWithPrefix(ctx, prefix)
		if err != nil {
			log.Errorf("error deleting objects with prefix %q", prefix)
			continue
		}
		log.Debugf("succesfully deleted %d objects with prefix %s", len(deleted), prefix)
		allDeleted = append(allDeleted, deleted...)
	}

	return allDeleted, nil
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
