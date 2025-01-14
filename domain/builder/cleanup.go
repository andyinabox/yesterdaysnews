package builder

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/streams"
)

func (b *Builder) Cleanup(ctx context.Context, toKeep int) ([]string, error) {

	deletedObjects := []string{}

	prefixes, err := b.cs.ListPrefixes(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing container prefixes: %w", err)
	}

	// first we find all the prefixes that are missing the manifest and therefore a failed build
	prefixesStream := streams.StringStreamThrottled(ctx, time.Millisecond, prefixes...)
	failedBuildPrefixes := b.failedBuildPrefixStream(ctx, prefixesStream)
	toDelete := streams.StringSlice(ctx, failedBuildPrefixes)

	// next we filter those out of our list
	remainingPrefixes := []string{}
	for _, pre := range prefixes {
		// skip any prefixes that are already in toDelete list
		for _, toDel := range toDelete {
			if toDel == pre {
				continue
			}
		}
		remainingPrefixes = append(remainingPrefixes, pre)
	}

	// finally only keep n of the remaining valid prefixes
	toDelete = append(toDelete, getPrefixesToDelete(remainingPrefixes, toKeep)...)

	if len(toDelete) != 0 {
		log.Infof("found %d prefixes to delete: %v", len(toDelete), toDelete)

		// convert prefix list to stream
		prefixesToDeleteStream := streams.StringStreamThrottled(ctx, time.Millisecond, toDelete...)

		// get object keys from prefixes
		objectsWithPrefixStream := b.cs.ListObjectsWithPrefixStream(ctx, b.errs, prefixesToDeleteStream)

		// add throttling
		objectsToDeleteStream := streams.StringPipeThrottled(ctx, time.Millisecond, objectsWithPrefixStream)

		// delete objects
		deletedStream := b.cs.DeleteObjectStream(ctx, b.errs, objectsToDeleteStream)

		// condense into slice
		deletedObjects = streams.StringSlice(ctx, deletedStream)
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

func (b *Builder) failedBuildPrefixStream(ctx context.Context, prefixes <-chan string) <-chan string {
	stream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		close(stream)
	}

	checkForManifest := func(pre string) {
		defer wg.Done()

		// we are assuming that if there is no manifest, the build failed
		key := fmt.Sprintf("%s/%s", pre, domain.ManifestFileName)
		log.Debugf("checking to see if %q exists", key)
		exists, err := b.cs.ObjectExists(ctx, key)
		if err != nil {
			b.errs <- errorhandler.Err(domain.ErrTypeCleanup, fmt.Errorf("error checking if object %q exists: %w", key, err))
			return
		}

		if !exists {
			log.Infof("prefix %q does not have a manifest, adding to delete list", pre)
			stream <- pre
		}
	}

	go func() {
		defer cleanup()
		for pre := range prefixes {
			select {
			case <-ctx.Done():
				return
			default:
				wg.Add(1)
				go checkForManifest(pre)
			}
		}
	}()

	return stream
}

func getPrefixesToDelete(prefixes []string, totalToKeep int) []string {

	if len(prefixes) <= totalToKeep {
		return []string{}
	}

	// sort prefixes by number
	sort.Slice(prefixes, func(i, j int) bool {
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
	})

	return prefixes[totalToKeep:]
}
