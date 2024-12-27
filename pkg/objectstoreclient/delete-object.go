package objectstoreclient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func (c *Client) DeleteObject(ctx context.Context, containerName string, key string) error {
	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	_, err = client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(containerName),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}
