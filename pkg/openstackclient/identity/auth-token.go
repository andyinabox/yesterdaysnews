package identity

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const authTokenRequestTemplate = `
{
	"auth": {
			"identity": {
					"methods": [
							"password"
					],
					"password": {
							"user": {
									"id": "%s",
									"password": "%s"
							}
					}
			},
			"scope": "unscoped"
	}
}	
`

func (c *Client) AuthToken(ctx context.Context) (string, error) {
	endpoint := fmt.Sprintf("%s/auth/tokens", c.cfg.IdentityEndpoint)
	reqBody := strings.TrimSpace(fmt.Sprintf(
		authTokenRequestTemplate,
		c.cfg.UserID,
		c.cfg.Password,
	))

	// log.Debug(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader([]byte(reqBody)))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("recieved a non-200 status code: %s", resp.Status)
	}

	token := resp.Header.Get("X-Subject-Token")

	if token == "" {
		return "", errors.New("no token recieved")
	}

	return token, nil
}
