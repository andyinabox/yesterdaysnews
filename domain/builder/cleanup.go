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
		// note that prefix includes "/" so we don't need to add here
		key := fmt.Sprintf("%s%s", pre, domain.ManifestFileName)

		log.Debugf("checking to see if %q exists", key)
		exists, err := b.cs.ObjectExists(ctx, key)
		if err != nil {
			b.errs <- errorhandler.Err(domain.ErrTypeCleanup, fmt.Errorf("error checking if object %q exists: %w", key, err))
			return true
		}
		log.Debugf("%q exists? %v", key, exists)

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
				log.Debugf("checking if %q should be deleted", prefix)

				// skip until we've reached our limit
				if kept < toKeep {
					log.Debugf("only %d prefixes have been kept, keeping %q", kept, prefix)
					kept++
					continue
				}

				// elete the rest
				log.Debugf("%d prefixes have already been kept, deleting %q", kept, prefix)
				toDelete <- prefix
			}
		}
	}()

	return toDelete
}
