package objectstoreclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *Client) ListPrefixes(ctx context.Context, containerName string) (prefixes []string, err error) {
	var result *s3.ListObjectsV2Output
	prefixesMap := make(map[types.CommonPrefix]struct{})

	client, err := c.getClient(ctx)
	if err != nil {
		return
	}

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(containerName),
		Delimiter: aws.String("/"),
	}

	paginator := s3.NewListObjectsV2Paginator(client, input)

	for paginator.HasMorePages() {

		result, err = paginator.NextPage(ctx)

		// handle error
		if err != nil {
			var noBucket *types.NoSuchBucket
			if errors.As(err, &noBucket) {
				err = fmt.Errorf("error listing objects: %w", ErrBucketDoesNotExist)
			}
			return
		}

		for _, prefix := range result.CommonPrefixes {
			prefixesMap[prefix] = struct{}{}
		}
	}

	prefixes = make([]string, len(prefixesMap))
	i := 0
	for prefix := range prefixesMap {
		prefixes[i] = *prefix.Prefix
		i++
	}

	return
}
