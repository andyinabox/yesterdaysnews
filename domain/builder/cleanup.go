package builder

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (b *Builder) Cleanup(ctx context.Context, toKeep string) ([]string, error) {

	var mu sync.Mutex

	log.Info("deleting prefixes except for %q and %q...", b.cfg.ObjectStorePrimaryDir, toKeep)
	deletedStream := b.deletePrefixesStream(ctx, toKeep)
	deletedObjects := []string{}

	for d := range deletedStream {
		mu.Lock()
		deletedObjects = append(deletedObjects, d)
		mu.Unlock()
	}

	log.Info("removing output dir")
	err := os.RemoveAll(b.cfg.OutputDir)
	if err != nil {
		return deletedObjects, fmt.Errorf("error removing %q from filesystem: %w", b.cfg.OutputDir, err)
	}

	return deletedObjects, nil
}

func (b *Builder) deletePrefixesStream(ctx context.Context, toKeep string) <-chan string {

	toDeleteStream := make(chan string)

	go func() {
		defer close(toDeleteStream)
		prefixes, err := b.cs.ListPrefixes(ctx)
		if err != nil {
			b.error(domain.ErrTypeTODO, fmt.Errorf("error listing container prefixes: %w", err))
		}

		for _, prefix := range prefixes {

			// skip primary dir
			if strings.HasPrefix(prefix, b.cfg.ObjectStorePrimaryDir) {
				continue
			}

			// skip dir to keep
			if toKeep != "" && strings.HasPrefix(prefix, toKeep) {
				continue
			}

			toDeleteStream <- prefix
		}

	}()

	return b.cs.DeleteObjectStream(ctx, b.errs, toDeleteStream)
}

// func (s *Service) PruneObjects(ctx context.Context, prefixToKeep string) ([]string, error) {

// 	prefixes, err := u.osclient.ListPrefixes(ctx, u.cfg.ContainerName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	allDeleted := []string{}

// 	for _, prefix := range prefixes {
// 		// skip primary dir
// 		if strings.HasPrefix(prefix, u.cfg.PrimaryDir) {
// 			continue
// 		}
// 		// skip keep prefix
// 		if prefixToKeep != "" && strings.HasPrefix(prefix, prefixToKeep) {
// 			continue
// 		}

// 		deleted, err := u.deleteObjectsWithPrefix(ctx, prefix)
// 		if err != nil {
// 			log.Errorf("error deleting objects with prefix %q", prefix)
// 			continue
// 		}
// 		log.Debugf("succesfully deleted %d objects with prefix %s", len(deleted), prefix)
// 		allDeleted = append(allDeleted, deleted...)
// 	}

// 	return allDeleted, nil
// }

// func (s *Service) deleteObjectsWithPrefix(ctx context.Context, prefix string) ([]string, error) {
// 	objectKeys, err := u.osclient.ListObjects(ctx, u.cfg.ContainerName, &objectstoreclient.ListObjectsRequest{
// 		Prefix: prefix,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("error listing objects: %w", err)
// 	}

// 	var wg sync.WaitGroup
// 	var mu sync.Mutex
// 	deleted := []string{}

// 	for _, key := range objectKeys {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()

// 			log.Debugf("deleting %q", key)
// 			err := u.osclient.DeleteObject(ctx, u.cfg.ContainerName, key)
// 			if err != nil {
// 				log.Errorf("error deleting %q: %s", key, err)
// 			}
// 			log.Debugf("finished deleting %q", key)

// 			mu.Lock()
// 			deleted = append(deleted, key)
// 			mu.Unlock()
// 		}()
// 	}

// 	wg.Wait()

// 	return deleted, nil
// }
