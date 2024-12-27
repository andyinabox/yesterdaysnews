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

type ListObjectsRequest struct {
	Prefix string
}

func (c *Client) ListObjects(ctx context.Context, containerName string, req *ListObjectsRequest) (fileKeys []string, err error) {
	var output *s3.ListObjectsV2Output
	var objects []types.Object

	client, err := c.getClient(ctx)
	if err != nil {
		return
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(containerName),
	}

	if req.Prefix != "" {
		input.Prefix = aws.String(req.Prefix)
	}

	paginator := s3.NewListObjectsV2Paginator(client, input)

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
