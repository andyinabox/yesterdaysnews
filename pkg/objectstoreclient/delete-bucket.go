package objectstoreclient

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
)

func (c *Client) DeleteBucket(ctx context.Context, bucketName string, force bool) error {
	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	objectKeys, err := c.ListObjects(ctx, bucketName)
	if err != nil {
		return err
	}

	// handle non-empty bucket
	if len(objectKeys) != 0 {

		// if force = true, delete bucket contents
		if force {
			log.Debugf("deleting %d objects from bucket %q prior to bucket deletion", len(objectKeys), bucketName)
			c.deleteBucketObjects(ctx, bucketName, objectKeys)

			// otherwise return an error
		} else {
			return fmt.Errorf("cannot delete non-empty bucket %q, use 'force=true'", bucketName)
		}
	}

	// this is hacky, but without it sometimes we get a
	// "bucket not empty" error when trying to delete
	time.Sleep(time.Second)

	log.Debugf("deleting bucket %q", bucketName)
	_, err = client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})

	return err
}

func (c *Client) deleteBucketObjects(ctx context.Context, bucketName string, keys []string) {
	var wg sync.WaitGroup
	// maxConcurrent := 20
	// totalConcurrent := 0

	for _, key := range keys {
		wg.Add(1)
		// totalConcurrent++

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

		// if totalConcurrent >= maxConcurrent {
		// 	wg.Wait()
		// 	totalConcurrent = 0
		// }
	}

	wg.Wait()
}
