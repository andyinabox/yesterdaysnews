package objectstoreclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
)

var ErrBucketDoesNotExist = fmt.Errorf("bucket does not exist")

func (c *Client) ListObjects(ctx context.Context, bucketName string) (fileKeys []string, err error) {
	var output *s3.ListObjectsV2Output
	var objects []types.Object

	client, err := c.getClient(ctx)
	if err != nil {
		return
	}

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})

	for paginator.HasMorePages() {

		output, err = paginator.NextPage(ctx)

		// handle error
		if err != nil {
			var noBucket *types.NoSuchBucket
			if errors.As(err, &noBucket) {
				err = fmt.Errorf("error listing objects: %w", ErrBucketDoesNotExist)
			}
			return
		}

		objects = append(objects, output.Contents...)
	}

	fileKeys = make([]string, len(objects))

	for i, o := range objects {
		fileKeys[i] = *o.Key
	}

	return
}
