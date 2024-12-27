package objectstoreclient

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Client) GetObject(ctx context.Context, containerName string, key string) ([]byte, error) {
	client, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(containerName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching object %q contents: %w", key, err)
	}

	defer result.Body.Close()
	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading object %q contents: %w", key, err)
	}

	return data, nil
}
