package objectstoreclient

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ratelimit"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
}

type Client struct {
	cfg        *Config
	awsconfig  *aws.Config
	s3Client   *s3.Client
	s3Uploader *manager.Uploader
}

func New(cfg *Config) *Client {

	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	return &Client{
		cfg: cfg,
	}
}

func (c *Client) getConfig(ctx context.Context) (*aws.Config, error) {
	if c.awsconfig == nil {
		// Load the Shared AWS Configuration (~/.aws/config)
		awsconfig, err := config.LoadDefaultConfig(
			ctx,
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				c.cfg.AccessKey,
				c.cfg.SecretKey,
				"",
			)),
			config.WithRegion(c.cfg.Region),
			// retryer should help when running into rate limiting
			config.WithRetryer(func() aws.Retryer {
				return retry.NewStandard(func(o *retry.StandardOptions) {
					// https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/retries-timeouts/#example-modified-rate-limiter
					// Don't totally understand what these values mean, but they seem to work better than the defaults
					o.RateLimiter = ratelimit.NewTokenRateLimit(1000)
					o.RetryCost = 1
					o.RetryTimeoutCost = 3
					o.NoRetryIncrement = 10
				})

			}),
		)
		if err != nil {
			return nil, fmt.Errorf("error loading config: %w", err)
		}
		c.awsconfig = &awsconfig
	}
	return c.awsconfig, nil
}

func (c *Client) getClient(ctx context.Context) (*s3.Client, error) {

	if c.s3Client == nil {
		awsconfig, err := c.getConfig(ctx)
		if err != nil {
			return nil, err
		}

		// Create an Amazon S3 service client
		client := s3.NewFromConfig(*awsconfig, func(o *s3.Options) {
			// force path style for openstack compatibility
			o.UsePathStyle = true
			o.BaseEndpoint = &c.cfg.Endpoint
		})
		c.s3Client = client
	}

	return c.s3Client, nil
}

func (c *Client) getUploader(ctx context.Context) (*manager.Uploader, error) {
	if c.s3Uploader == nil {
		client, err := c.getClient(ctx)
		if err != nil {
			return nil, err
		}

		uploader := manager.NewUploader(client)
		c.s3Uploader = uploader
	}

	return c.s3Uploader, nil
}
