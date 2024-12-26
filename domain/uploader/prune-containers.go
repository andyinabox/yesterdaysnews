package uploader

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/objectstoreclient"
)

func (u *Uploader) PruneContainers(ctx context.Context, containerToKeep string) ([]string, error) {

	containers, err := u.osclient.ListBuckets(ctx, &objectstoreclient.ListBucketsRequest{
		Prefix: fmt.Sprintf("%s-", u.cfg.BucketNameBase),
	})
	if err != nil {
		return nil, err
	}

	var removed []string

	for _, containerName := range containers {
		if containerName != containerToKeep {
			log.Debugf("deleting container %q", containerName)
			err = u.osclient.DeleteBucket(ctx, containerName, true)
			if err != nil {
				return removed, fmt.Errorf("error deleting container %q: %w", containerName, err)
			}
			removed = append(removed, containerName)
		}
	}

	return removed, nil
}
