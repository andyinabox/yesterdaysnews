package objectstoreclient

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
)

func (c *Client) DeleteBucket(ctx context.Context, bucketName string, force bool) error {
	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	// we are looping through the list command because sometimes
	// the delete command is failing with a "bucket not empty" error
	for {
		objectKeys, err := c.ListObjects(ctx, bucketName)
		if err != nil {
			return err
		}

		// break out of loop when there are no more objects left
		if len(objectKeys) == 0 {
			break
		}

		// if force is not true, return error
		if !force {
			return fmt.Errorf("cannot delete non-empty bucket %q, use 'force=true'", bucketName)
		}

		log.Debugf("deleting %d objects from bucket %q prior to bucket deletion", len(objectKeys), bucketName)
		c.deleteBucketObjects(ctx, bucketName, objectKeys)
	}

	log.Debugf("deleting bucket %q", bucketName)
	_, err = client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})

	return err
}

func (c *Client) deleteBucketObjects(ctx context.Context, bucketName string, keys []string) {
	var wg sync.WaitGroup

	for _, key := range keys {
		wg.Add(1)

		go func() {
			defer wg.Done()
			log.Debugf("deleting object %q from bucket %q", key, bucketName)
			err := c.DeleteObject(ctx, bucketName, key)
			if err != nil {
				log.Errorf("error deleting object %q from bucket %q: %s", key, bucketName, err)
				return
			}
			log.Debugf("successfully deleted %q", key)
		}()

	}

	wg.Wait()
}
