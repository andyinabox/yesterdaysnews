package osobjectstoreclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
)

type ContainerSummary struct {
	Count        int       `json:"count"`
	Bytes        int       `json:"bytes"`
	Name         string    `json:"name"`
	LastModified time.Time `json:"last_modified"`
}

type ListContainerResult struct {
	AccountObjectCount    int
	AccountBytesUsed      int
	AccountContainerCount int
	Containers            []ContainerSummary
}

func (c *Client) ListContainers(ctx context.Context) (*ListContainerResult, error) {
	endpoint := fmt.Sprintf("%s/%s", c.cfg.Endpoint, c.cfg.ProjectID)

	log.Debug("creating request", "endpoint", endpoint)

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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	containers := []ContainerSummary{}
	err = json.Unmarshal(body, &containers)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %w", err)
	}

	result := &ListContainerResult{
		Containers: containers,
	}

	// parse headers
	result.AccountObjectCount, err = strconv.Atoi(resp.Header.Get("X-Account-Object-Count"))
	if err != nil {
		log.Error("error parsing int AccountObjectCount")
	}
	result.AccountBytesUsed, err = strconv.Atoi(resp.Header.Get("X-Account-Bytes-Used"))
	if err != nil {
		log.Error("error parsing int AccountBytesUsed")
	}
	result.AccountContainerCount, err = strconv.Atoi(resp.Header.Get("X-Account-Container-Count"))
	if err != nil {
		log.Error("error parsing int AccountContainerCount")
	}

	return result, nil
}
