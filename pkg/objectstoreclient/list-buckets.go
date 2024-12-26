package objectstoreclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ListBucketsRequest struct {
	Prefix string
}

func (c *Client) ListBuckets(ctx context.Context, req *ListBucketsRequest) (buckets []string, err error) {
	client, err := c.getClient(ctx)
	if err != nil {
		err = fmt.Errorf("error getting client: %w", err)
		return
	}

	input := &s3.ListBucketsInput{}

	if req.Prefix != "" {
		input.Prefix = aws.String(req.Prefix)
	}

	result, err := client.ListBuckets(ctx, input)
	if err != nil {
		err = fmt.Errorf("error listing buckets: %w", err)
		return
	}

	buckets = []string{}
	for _, b := range result.Buckets {
		name := *b.Name
		// the `Prefix` option does not seem to be working for some reason so doing this manually
		if strings.HasPrefix(name, req.Prefix) {
			buckets = append(buckets, name)
		}
	}

	return
}
