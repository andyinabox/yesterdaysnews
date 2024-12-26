package objectstoreclient

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Client) GetObject(ctx context.Context, bucketName string, objectKey string) ([]byte, error) {
	client, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching object %q contents: %w", objectKey, err)
	}

	defer result.Body.Close()
	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading object %q contents: %w", objectKey, err)
	}

	return data, nil
}
