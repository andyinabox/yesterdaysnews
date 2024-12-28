package osobjectstoreclient

import "context"

type OSIdentityClient interface {
	AuthToken(context.Context) (string, error)
}

type Config struct {
	Endpoint   string
	ProjectID  string
	RegionName string
}

type Client struct {
	auth OSIdentityClient
	cfg  *Config
}

func New(auth OSIdentityClient, cfg *Config) *Client {
	return &Client{auth, cfg}
}

func (c *Client) getToken(ctx context.Context) (string, error) {
	// TODO: cache token?
	return c.auth.AuthToken(ctx)
}
