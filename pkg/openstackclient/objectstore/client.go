package objectstore

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
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

func (c *Client) doPostHeadersRequest(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.cfg.Endpoint, path)

	log.Debugf("executing request for %q", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range headers {
		log.Debugf("set header %q: %q", key, value)
		req.Header.Set(key, value)
	}

	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting token: %w", err)
	}

	log.Debug("setting auth header", "token", token)
	req.Header.Set("X-Auth-Token", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("recieved a non-200 status code: %s", resp.Status)
	}

	return resp, nil
}

func (c *Client) doGetRequest(ctx context.Context, path string, query url.Values) ([]byte, *http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.cfg.Endpoint, path)

	if query != nil {
		url = url + "?" + query.Encode()
	}

	log.Debugf("executing request for %q", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %w", err)
	}

	token, err := c.getToken(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting token: %w", err)
	}

	log.Debug("setting auth header", "token", token)
	req.Header.Set("X-Auth-Token", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("recieved a non-200 status code: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, resp, nil
}
