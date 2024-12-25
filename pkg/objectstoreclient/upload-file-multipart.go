package objectstoreclient

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Client) UploadFileMultipart(ctx context.Context, bucketName, fileKey, inFilepath string) (string, error) {

	file, err := os.Open(inFilepath)
	if err != nil {
		return "", fmt.Errorf("unable to open file %q: %w", inFilepath, err)
	}

	uploader, err := c.getUploader(ctx)
	if err != nil {
		return "", err
	}

	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &fileKey,
		Body:   file,
	})

	return fileKey, nil
}
