package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func (b *Builder) Promote(ctx context.Context, uploadDir string) error {

	// first move the files from the primary prefix to a different prefiex
	demotedName := b.getCurrentManifestName(ctx)
	log.Infof("demoting current objects to %q", demotedName)
	fileKeys, err := b.moveObjects(ctx, b.cfg.ObjectStorePrimaryDir, demotedName)
	if err != nil {
		return fmt.Errorf("error moving objects from %q to %q: %w", b.cfg.ObjectStorePrimaryDir, demotedName, err)
	}
	log.Infof("successfully moved %d objects from %q to %q", len(fileKeys), b.cfg.ObjectStorePrimaryDir, demotedName)

	log.Info("promoting %q", uploadDir)
	fileKeys, err = b.moveObjects(ctx, uploadDir, b.cfg.ObjectStorePrimaryDir)
	if err != nil {
		return fmt.Errorf("error moving objects from %q to %q: %w", uploadDir, b.cfg.ObjectStorePrimaryDir, err)
	}
	log.Infof("successfully moved %d objects from %q to %q", len(fileKeys), uploadDir, b.cfg.ObjectStorePrimaryDir)

	return nil
}

func (b *Builder) moveObjects(ctx context.Context, fromPrefix, toPrefix string) ([]string, error) {
	stream := make(chan [2]string)

	var mu sync.Mutex

	movedFiles := []string{}

	fileKeys, err := b.os.ListObjectsInDir(ctx, fromPrefix)
	if err != nil {
		return nil, fmt.Errorf("error listing objects with %q prefix: %w", fromPrefix, err)
	}

	// tranlate to use correct format for object move stream
	go func() {
		defer close(stream)
		for _, fromKey := range fileKeys {
			toKey := strings.Replace(fromKey, fromPrefix, toPrefix, 1)
			stream <- [2]string{fromKey, toKey}
		}
	}()

	for key := range b.os.MoveObjectStream(ctx, b.errs, stream) {
		mu.Lock()
		defer mu.Unlock()
		movedFiles = append(movedFiles, key)
	}

	return movedFiles, nil
}

// getCurrentManifestName attempts to get the build ID from the manifest in
// the primary object store prefix/dir. If that fails default to a current timestamp
func (b *Builder) getCurrentManifestName(ctx context.Context) (prefix string) {
	var err error
	prefix = util.Timestamp(time.Now())

	currentManifestFileKey := filepath.Join(b.cfg.ObjectStorePrimaryDir, "manifest.json")
	exists, err := b.os.ObjectExists(ctx, currentManifestFileKey)
	if err != nil {
		b.error(domain.ErrTypeGetCurrentManifestPrefix, fmt.Errorf("error determining if %q exists: %w", err))
		return
	}

	// attempt to get the
	if exists {
		data, err := b.os.DownloadObject(ctx, currentManifestFileKey)
		if err != nil {
			b.error(domain.ErrTypeGetCurrentManifestPrefix, fmt.Errorf("error downloading %q: %w", currentManifestFileKey, err))
			return
		}

		manifest := domain.Manifest{}
		err = json.Unmarshal(data, &manifest)
		if err != nil {
			b.error(domain.ErrTypeGetCurrentManifestPrefix, fmt.Errorf("error unmarshaling current manifest: %w", err))
			return
		}

		if manifest.ID != "" {
			prefix = manifest.ID
		}

	}

	return
}
