package objectstoreclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

func (c *Client) ContainerExists(ctx context.Context, name string) (bool, error) {
	client, err := c.getClient(ctx)
	if err != nil {
		return false, err
	}

	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &name,
	})

	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				return false, nil
			default:
				return false, fmt.Errorf("error checking if bucket exists: %w", err)
			}
		}
	}

	return true, nil
}
