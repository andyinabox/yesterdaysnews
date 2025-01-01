package uploader

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/objectstoreclient"
)

func (u *Uploader) PromoteObjects(ctx context.Context, prefix string) (string, error) {

	currentObjects, err := u.osclient.ListObjects(ctx, u.cfg.ContainerName, &objectstoreclient.ListObjectsRequest{
		Prefix: u.cfg.PrimaryDir,
	})
	if err != nil {
		return "", err
	}

	var demotedDir string

	if len(currentObjects) != 0 {
		demotedDir = u.getDemotedDirName(ctx)
		err = u.moveObjects(ctx, u.cfg.PrimaryDir, demotedDir)
		if err != nil {
			return "", fmt.Errorf("error demoting current object: %w", err)
		}
	}

	return demotedDir, u.moveObjects(ctx, prefix, u.cfg.PrimaryDir)
}

func (u *Uploader) getDemotedDirName(ctx context.Context) (dir string) {
	dir = timestamp()

	data, err := u.osclient.GetObject(ctx, u.cfg.ContainerName, filepath.Join(u.cfg.PrimaryDir, "manifest.json"))
	if err != nil {
		log.Errorf("error fetching current manifest: %s", err)
		return
	}

	manifest := domain.Manifest{}
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		log.Errorf("error unmarshaling current manifest: %s", err)
		return
	}

	if manifest.ID != "" {
		dir = manifest.ID
	}

	return
}

func (u *Uploader) moveObjects(ctx context.Context, sourcePrefix, destPrefix string) error {
	toMove, err := u.osclient.ListObjects(ctx, u.cfg.ContainerName, &objectstoreclient.ListObjectsRequest{
		Prefix: sourcePrefix,
	})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, key := range toMove {

		wg.Add(1)
		go func() {
			defer wg.Done()
			source := key
			dest := strings.Replace(key, sourcePrefix, destPrefix, 1)

			log.Debugf("copying %q to %q", source, dest)
			err := u.osclient.CopyObject(
				ctx,
				u.cfg.ContainerName,
				u.cfg.ContainerName,
				source,
				dest,
			)
			if err != nil {
				log.Errorf("error moving %q to %q", source, dest)
				return
			}
			err = u.osclient.DeleteObject(ctx, u.cfg.ContainerName, source)
			if err != nil {
				log.Errorf("error deleting object %q", source)
				return
			}
		}()

	}

	wg.Wait()

	return nil
}
