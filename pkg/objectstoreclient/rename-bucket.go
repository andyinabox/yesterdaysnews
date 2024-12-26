package objectstoreclient

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
)

var ErrBucketAlreadyExists = errors.New("bucket already exists")

func (c *Client) RenameBucket(ctx context.Context, currentName, newName string) error {

	exists, err := c.BucketExists(ctx, newName)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("error renaming %q to %q: %w", currentName, newName, ErrBucketAlreadyExists)
	}

	err = c.CreateBucket(ctx, newName)
	if err != nil {
		return fmt.Errorf("error creating bucket %q: %w", newName, err)
	}

	keys, err := c.ListObjects(ctx, currentName)
	if err != nil {
		return fmt.Errorf("error listing objects from bucket %q: %w", currentName, err)
	}

	var wg sync.WaitGroup

	for _, key := range keys {

		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Debugf("copy object %q from %q to %q", key, currentName, newName)
			err = c.CopyObjectToBucket(ctx, currentName, newName, key)
			if err != nil {
				log.Errorf("error moving %q to bucket %q: %s", key, newName, err)
				return
			}
			log.Debugf("finished copying object %q", key)
		}()

	}

	wg.Wait()

	err = c.DeleteBucket(ctx, currentName, true)
	if err != nil {
		return fmt.Errorf("error deleting old bucket %q: %w", currentName, err)
	}

	return nil
}
