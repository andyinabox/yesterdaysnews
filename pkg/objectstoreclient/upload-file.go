package objectstoreclient

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Client) UploadFile(ctx context.Context, bucketName, fileKey, inFilepath string) (string, error) {
	file, err := os.Open(inFilepath)
	if err != nil {
		return "", fmt.Errorf("unable to open file %q: %w", inFilepath, err)
	}

	client, err := c.getClient(ctx)
	if err != nil {
		return "", nil
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &fileKey,
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("error putting file %q: %w", fileKey, err)
	}

	return fileKey, nil
}
