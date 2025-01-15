package objectstoreclient

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *Client) UploadPublicFile(ctx context.Context, containerName, fileKey string, reader io.Reader, contentType string, multipart bool) (string, error) {
	if multipart {
		return c.uploadPublicFileMultipart(ctx, containerName, fileKey, reader, contentType)
	}

	return c.uploadPublicFile(ctx, containerName, fileKey, reader, contentType)
}
func (c *Client) uploadPublicFile(ctx context.Context, containerName, fileKey string, reader io.Reader, contentType string) (string, error) {

	client, err := c.getClient(ctx)
	if err != nil {
		return "", nil
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(containerName),
		Key:         aws.String(fileKey),
		Body:        reader,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("error putting file %q: %w", fileKey, err)
	}

	return fileKey, nil
}

func (c *Client) uploadPublicFileMultipart(ctx context.Context, containerName, fileKey string, reader io.Reader, contentType string) (string, error) {

	uploader, err := c.getUploader(ctx)
	if err != nil {
		return "", err
	}

	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(containerName),
		Key:         aws.String(fileKey),
		Body:        reader,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("error uploading file %q: %w", fileKey, err)
	}

	return fileKey, nil
}
