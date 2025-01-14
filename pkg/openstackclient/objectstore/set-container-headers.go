package objectstore

import "context"

func (c *Client) SetContainerMetadata(ctx context.Context, containerName string, headers map[string]string) error {
	_, err := c.doPostHeadersRequest(ctx, containerName, headers)
	return err
}
