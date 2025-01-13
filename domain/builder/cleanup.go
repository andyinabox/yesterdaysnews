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
	"gitlab.com/andyinabox/yesterdaysnews/pkg/streams"
)

func (b *Builder) Cleanup(ctx context.Context, toKeep int) ([]string, error) {

	deletedObjects := []string{}

	prefixes, err := b.cs.ListPrefixes(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing container prefixes: %w", err)
	}

	toDelete := getPrefixesToDelete(prefixes, toKeep)

	if len(toDelete) != 0 {
		log.Infof("found %d prefixes to delete: %v", len(toDelete), toDelete)
		prefixesStream := streams.StringStreamThrottled(ctx, time.Millisecond, toDelete...)

		// filteredPrefixesStream := streams.StringFilterStream(ctx, prefixesStream, func(s string) bool {
		// 	log.Debugf("checking %q", s)

		// 	// skip primary dir
		// 	if strings.HasPrefix(s, b.cfg.ObjectStorePrimaryDir) {
		// 		return false
		// 	}

		// 	// skip dir to keep
		// 	if toKeep != "" && strings.HasPrefix(s, toKeep) {
		// 		return false
		// 	}

		// 	log.Debugf("adding %q to delete stream", s)
		// 	return true
		// })

		objectsWithPrefixStream := b.cs.ListObjectsWithPrefixStream(ctx, b.errs, prefixesStream)

		// add throttling
		objectsToDeleteStream := streams.StringPipeThrottled(ctx, time.Millisecond, objectsWithPrefixStream)

		deletedStream := b.cs.DeleteObjectStream(ctx, b.errs, objectsToDeleteStream)

		deletedObjects = streams.StringSlice(ctx, deletedStream)
	}

	if b.cfg.RemoveFilesOnCompletion {
		log.Info("removing output dir")
		err = os.RemoveAll(b.cfg.OutputDir)
		if err != nil {
			return deletedObjects, fmt.Errorf("error removing %q from filesystem: %w", b.cfg.OutputDir, err)
		}
	}

	return deletedObjects, nil
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
