package objectstoreclient

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func (c *Client) CopyObjectToBucket(ctx context.Context, sourceBucket, destBucket string, fileKey string) error {
	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(destBucket),
		CopySource: aws.String(fmt.Sprintf("%s/%s", sourceBucket, fileKey)),
		Key:        aws.String(fileKey),
	})

	return err
}
