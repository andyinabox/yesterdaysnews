package objectstoreclient

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Client) UploadFileMultipart(ctx context.Context, bucketName, fileKey string, reader io.Reader, contentType string) (string, error) {

	uploader, err := c.getUploader(ctx)
	if err != nil {
		return "", err
	}

	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileKey),
		Body:        reader,
		ContentType: aws.String(contentType),
	})

	return fileKey, nil
}
