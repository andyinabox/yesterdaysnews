package objectstoreclient

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
}

type Client struct {
	cfg    *Config
	client *s3.Client
}

func New(cfg *Config) *Client {

	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	return &Client{
		cfg: cfg,
	}
}

func (c *Client) getClient(ctx context.Context) (*s3.Client, error) {

	if c.client == nil {
		// Load the Shared AWS Configuration (~/.aws/config)
		cfg, err := config.LoadDefaultConfig(
			ctx,
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				c.cfg.AccessKey,
				c.cfg.SecretKey,
				"",
			)),
			config.WithRegion(c.cfg.Region),
		)
		if err != nil {
			return nil, fmt.Errorf("error loading config: %w", err)
		}

		// Create an Amazon S3 service client
		client := s3.NewFromConfig(cfg, func(o *s3.Options) {
			// force path style for openstack compatibility
			o.UsePathStyle = true
			o.BaseEndpoint = &c.cfg.Endpoint
		})
		c.client = client
	}

	return c.client, nil
}
