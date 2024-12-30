package objectstore

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
)

func (c *Client) Info(ctx context.Context) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/info", c.cfg.Endpoint)

	log.Debugf("endpoint: %s", endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
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

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
