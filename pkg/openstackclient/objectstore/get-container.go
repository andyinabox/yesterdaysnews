package objectstore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/charmbracelet/log"
)

type GetContainerResponse struct {
	ObjectCount int
	BytesUsed   int
	Objects     []ContainerObject
}

type ContainerObject struct {
	Hash         string `json:"hash"`
	LastModified string `json:"last_modified"`
	Bytes        uint   `json:"bytes"`
	Name         string `json:"name"`
	ContentType  string `json:"content_type"`
}

func (c *Client) GetContainer(ctx context.Context, containerName string) (*GetContainerResponse, error) {

	query := url.Values{}
	query.Set("format", "json")

	body, resp, err := c.doGetRequest(ctx, containerName, query)
	if err != nil {
		return nil, err
	}

	result := GetContainerResponse{
		Objects: []ContainerObject{},
	}

	err = json.Unmarshal(body, &result.Objects)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling GetContainer response: %w", err)
	}

	result.ObjectCount, err = strconv.Atoi(resp.Header.Get("X-Container-Object-Count"))
	if err != nil {
		log.Error("error parsing int ObjectCount")
	}

	result.BytesUsed, err = strconv.Atoi(resp.Header.Get("X-Container-Bytes-Used"))
	if err != nil {
		log.Error("error parsing int BytesUsed")
	}

	return &result, nil
}
