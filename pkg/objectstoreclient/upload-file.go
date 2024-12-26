package objectstoreclient

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Client) UploadFile(ctx context.Context, bucketName, fileKey string, reader io.Reader, contentType string) (string, error) {

	client, err := c.getClient(ctx)
	if err != nil {
		return "", nil
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileKey),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("error putting file %q: %w", fileKey, err)
	}

	return fileKey, nil
}
