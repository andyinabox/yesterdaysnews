package objectstoreclient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *Client) CreateBucket(ctx context.Context, name string) error {
	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &name,
		ACL:    types.BucketCannedACLPublicRead,
	})
	if err != nil {
		return err
	}

	return nil
}
