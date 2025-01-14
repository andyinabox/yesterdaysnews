package builder

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/streams"
)

func (b *Builder) Cleanup(ctx context.Context, toKeep int) ([]string, error) {

	deletedObjects, err := b.cleanupObjectStore(ctx, toKeep)
	if err != nil {
		return nil, fmt.Errorf("error cleaning up object store: %w", err)
	}

	if b.cfg.RemoveFilesOnCompletion {
		log.Info("removing dir %q", b.cfg.OutputDir)
		err = os.RemoveAll(b.cfg.OutputDir)
		if err != nil {
			return deletedObjects, fmt.Errorf("error removing %q from filesystem: %w", b.cfg.OutputDir, err)
		}
	}

	return deletedObjects, nil
}

func (b *Builder) cleanupObjectStore(ctx context.Context, toKeep int) ([]string, error) {

	prefixes, err := b.cs.ListPrefixes(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing container prefixes: %w", err)
	}

	// sort prefix timestamps in descending order
	prefixeSortFn := func(i, j int) bool {
		s1 := strings.ReplaceAll(prefixes[i], "/", "")
		n1, err := strconv.Atoi(s1)
		if err != nil {
			return false
		}
		s2 := strings.ReplaceAll(prefixes[j], "/", "")
		n2, err := strconv.Atoi(s2)
		if err != nil {
			return true
		}
		return n1 > n2
	}

	// returns false if manifest.json does not exist, on the assumption that
	// it means the build did not finished and can be deleted
	prefixesBifurcateFn := func(pre string) bool {
		key := fmt.Sprintf("%s/%s", pre, domain.ManifestFileName)
		log.Debugf("checking to see if %q exists", key)
		exists, err := b.cs.ObjectExists(ctx, key)
		if err != nil {
			b.errs <- errorhandler.Err(domain.ErrTypeCleanup, fmt.Errorf("error checking if object %q exists: %w", key, err))
			return true
		}

		return exists
	}

	// sort prefixes descending
	sort.Slice(prefixes, prefixeSortFn)
	// convert into throttled stream
	prefixesStream := streams.StringStreamThrottled(ctx, time.Millisecond, prefixes...)
	// fork into prefixes to be culled, and prefixes from failed builds to be deleted
	toCull, toDeleteFailedBuild := streams.BifurcatedStringStream(ctx, prefixesStream, prefixesBifurcateFn)
	// cull down to selected number to keep, deleting the rest
	toDeleteFromCulling := b.getCullingToDelete(ctx, toCull, toKeep)
	// merge all prefixes to be deleted into single stream
	prefixesToDelete := streams.MergeStringStreams(ctx, toDeleteFailedBuild, toDeleteFromCulling)
	// get objects for each prefix so they can be deleted
	objectsToDelete := b.cs.ListObjectsWithPrefixStream(ctx, b.errs, prefixesToDelete)
	// throttle objects to hopefully reduce "To Many Requests" errors
	throttledObjectsToDelete := streams.StringPipeThrottled(ctx, time.Millisecond, objectsToDelete)
	// delete objects
	deletedObjectsStream := b.cs.DeleteObjectStream(ctx, b.errs, throttledObjectsToDelete)

	return streams.StringSlice(ctx, deletedObjectsStream), nil
}

func (b *Builder) getCullingToDelete(ctx context.Context, toCull <-chan string, toKeep int) <-chan string {
	toDelete := make(chan string)

	var kept int

	go func() {
		defer close(toDelete)
		for prefix := range toCull {
			select {
			case <-ctx.Done():
				return
			default:
				// skip until we've reached our limit
				if kept < toKeep {
					kept++
					continue
				}

				// elete the rest
				toDelete <- prefix
			}
		}
	}()

	return toDelete
}

// func (b *Builder) cleanupFailedBuilds(ctx context.Context, prefixes <-chan string) (<-chan string, <-chan string) {
// 	toDelete := make(chan string)
// 	remainder := make(chan string)

// 	var wg sync.WaitGroup

// 	cleanup := func() {
// 		wg.Wait()
// 		close(toDelete)
// 		close(remainder)
// 	}

// 	checkForManifest := func(pre string) {
// 		defer wg.Done()
// 		key := fmt.Sprintf("%s/%s", pre, domain.ManifestFileName)
// 		log.Debugf("checking to see if %q exists", key)
// 		exists, err := b.cs.ObjectExists(ctx, key)
// 		if err != nil {
// 			b.errs <- errorhandler.Err(domain.ErrTypeCleanup, fmt.Errorf("error checking if object %q exists: %w", key, err))
// 			remainder <- pre
// 			return
// 		}

// 		if exists {
// 			remainder <- pre
// 			return
// 		}

// 		toDelete <- pre
// 	}

// 	go func() {
// 		defer cleanup()
// 		for pre := range prefixes {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			default:
// 				wg.Add(1)
// 				go checkForManifest(pre)
// 			}
// 		}
// 	}()

// 	return toDelete, remainder
// }

// func (b *Builder) cleanupObjectStore(ctx context.Context, toKeep int) ([]string, error) {

// 	prefixes, err := b.cs.ListPrefixes(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("error listing container prefixes: %w", err)
// 	}

// 	deleteStream, closeDeleteStream := b.deletePrefixStream(ctx)

// 	// sort prefixes descending
// 	sort.Slice(prefixes, func(i, j int) bool {
// 		s1 := strings.ReplaceAll(prefixes[i], "/", "")
// 		n1, err := strconv.Atoi(s1)
// 		if err != nil {
// 			return false
// 		}
// 		s2 := strings.ReplaceAll(prefixes[j], "/", "")
// 		n2, err := strconv.Atoi(s2)
// 		if err != nil {
// 			return true
// 		}
// 		return n1 > n2
// 	})

// 	var totalKept int
// 	for _, pre := range prefixes {
// 		if totalKept < toKeep {

// 			// check for manifest
// 			key := fmt.Sprintf("%s/%s", pre, domain.ManifestFileName)
// 			log.Debugf("checking to see if %q exists", key)
// 			exists, err := b.cs.ObjectExists(ctx, key)
// 			if err != nil {
// 				b.errs <- errorhandler.Err(domain.ErrTypeCleanup, fmt.Errorf("error checking if object %q exists: %w", key, err))
// 				continue
// 			}

// 			// keep this one
// 			if exists {
// 				log.Debugf("keeping %q", pre)
// 				totalKept++
// 				continue
// 			}
// 		}

// 		log.Debugf("deleting prefix %q", pre)
// 		deleteStream <- pre
// 	}

// 	return closeDeleteStream(), nil

// }

// func (b *Builder) deletePrefixStream(ctx context.Context) (chan<- string, func() []string) {
// 	deletePrefixStream := make(chan string)
// 	done := make(chan struct{})

// 	allDeleted := []string{}

// 	var mu sync.Mutex
// 	var wg sync.WaitGroup

// 	doneFunc := func() []string {
// 		done <- struct{}{}
// 		wg.Wait()
// 		return allDeleted
// 	}

// 	cleanup := func() {
// 		wg.Wait()
// 		close(deletePrefixStream)
// 		close(done)
// 	}

// 	deleteObjectsForPrefix := func(prefix string) {
// 		defer wg.Done()

// 		log.Debugf("listing objects with prefix %q", prefix)
// 		objects, err := b.cs.ListObjectsWithPrefix(ctx, prefix)
// 		if err != nil {
// 			b.errs <- errorhandler.Err(domain.ErrTypeCleanup, fmt.Errorf("could not get objects for prefix %q: %w", prefix, err))
// 			return
// 		}
// 		log.Debugf("found %d objects with prefix %q", len(objects), prefix)

// 		objectsStream := streams.StringStreamThrottled(ctx, time.Millisecond, objects...)
// 		deletedStream := b.cs.DeleteObjectStream(ctx, b.errs, objectsStream)
// 		deleted := streams.StringSlice(ctx, deletedStream)
// 		log.Debugf("deleted %d objects with prefix %q", len(deleted), prefix)

// 		mu.Lock()
// 		allDeleted = append(allDeleted, deleted...)
// 		mu.Unlock()
// 	}

// 	go func() {
// 		defer cleanup()
// 		for prefix := range deletePrefixStream {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case <-done:
// 				return
// 			default:
// 				wg.Add(1)
// 				go deleteObjectsForPrefix(prefix)
// 			}
// 		}
// 	}()

// 	return deletePrefixStream, doneFunc
// }

// func (b *Builder) failedBuildPrefixStream(ctx context.Context, prefixes <-chan string, deleteStream chan<- string) <-chan string {
// 	stream := make(chan string)

// 	var wg sync.WaitGroup

// 	cleanup := func() {
// 		wg.Wait()
// 		close(stream)
// 	}

// 	checkForManifest := func(pre string) {
// 		defer wg.Done()

// 		// we are assuming that if there is no manifest, the build failed
// 		key := fmt.Sprintf("%s/%s", pre, domain.ManifestFileName)
// 		log.Debugf("checking to see if %q exists", key)
// 		exists, err := b.cs.ObjectExists(ctx, key)
// 		if err != nil {
// 			b.errs <- errorhandler.Err(domain.ErrTypeCleanup, fmt.Errorf("error checking if object %q exists: %w", key, err))
// 			return
// 		}

// 		if !exists {
// 			log.Infof("prefix %q does not have a manifest, adding to delete list", pre)
// 			stream <- pre
// 		}
// 	}

// 	go func() {
// 		defer cleanup()
// 		for pre := range prefixes {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			default:
// 				wg.Add(1)
// 				go checkForManifest(pre)
// 			}
// 		}
// 	}()

// 	return stream
// }

// func getPrefixesToDelete(prefixes []string, totalToKeep int) []string {

// 	if len(prefixes) <= totalToKeep {
// 		return []string{}
// 	}

// 	// sort prefixes by number
// 	sort.Slice(prefixes, func(i, j int) bool {
// 		s1 := strings.ReplaceAll(prefixes[i], "/", "")
// 		n1, err := strconv.Atoi(s1)
// 		if err != nil {
// 			return false
// 		}
// 		s2 := strings.ReplaceAll(prefixes[j], "/", "")
// 		n2, err := strconv.Atoi(s2)
// 		if err != nil {
// 			return true
// 		}
// 		return n1 > n2
// 	})

// 	return prefixes[totalToKeep:]
// }
