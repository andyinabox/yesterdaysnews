package objectstoreclient

import (
	"context"
	"fmt"
)

func (c *Client) ListBuckets(ctx context.Context) (buckets []string, err error) {
	client, err := c.getClient(ctx)
	if err != nil {
		err = fmt.Errorf("error getting client: %w", err)
		return
	}

	result, err := client.ListBuckets(ctx, nil)
	if err != nil {
		err = fmt.Errorf("error listing buckets: %w", err)
		return
	}

	buckets = make([]string, len(result.Buckets))
	for i, b := range result.Buckets {
		buckets[i] = *b.Name
	}

	return
}
