package objectstore

import (
	"context"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/openstackclient"
)

type Config struct {
	Endpoint   string
	ProjectID  string
	RegionName string
}

type Client struct {
	auth openstackclient.IdentityClient
	cfg  *Config
}

func New(auth openstackclient.IdentityClient, cfg *Config) *Client {
	return &Client{auth, cfg}
}

func (c *Client) getToken(ctx context.Context) (string, error) {
	// TODO: cache token?
	return c.auth.AuthToken(ctx)
}
