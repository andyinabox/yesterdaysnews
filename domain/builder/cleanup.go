package builder

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/streams"
)

func (b *Builder) Cleanup(ctx context.Context, toKeep string) ([]string, error) {

	prefixes, err := b.cs.ListPrefixes(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing container prefixes: %w", err)
	}

	prefixesStream := streams.StringStreamThrottled(ctx, time.Millisecond, prefixes...)

	filteredPrefixesStream := streams.StringFilterStream(ctx, prefixesStream, func(s string) bool {
		log.Debugf("checking %q", s)

		// skip primary dir
		if strings.HasPrefix(s, b.cfg.ObjectStorePrimaryDir) {
			return false
		}

		// skip dir to keep
		if toKeep != "" && strings.HasPrefix(s, toKeep) {
			return false
		}

		log.Debugf("adding %q to delete stream", s)
		return true
	})

	objectsWithPrefixStream := b.cs.ListObjectsWithPrefixStream(ctx, b.errs, filteredPrefixesStream)

	// add throttling
	objectsToDeleteStream := streams.StringPipeThrottled(ctx, time.Millisecond, objectsWithPrefixStream)

	deletedStream := b.cs.DeleteObjectStream(ctx, b.errs, objectsToDeleteStream)

	deletedObjects := streams.StringSlice(ctx, deletedStream)

	log.Info("removing output dir")
	err = os.RemoveAll(b.cfg.OutputDir)
	if err != nil {
		return deletedObjects, fmt.Errorf("error removing %q from filesystem: %w", b.cfg.OutputDir, err)
	}

	return deletedObjects, nil
}
