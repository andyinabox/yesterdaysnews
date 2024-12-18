package youtubeapi

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/youtubeapi/response"
)

const playlistItemsListEndpoint = "playlistItems"

type PlaylistItemsListRequest struct {
	PlaylistId string
	Part       []string
	MaxResults int
	PageToken  string
}

func (r *PlaylistItemsListRequest) Values() (q url.Values) {
	q = make(url.Values)

	q.Set("part", strings.Join(r.Part, ","))
	q.Set("playlistId", r.PlaylistId)
	q.Set("maxResults", strconv.Itoa(r.MaxResults)) // this is the max allowed

	if r.PageToken != "" {
		q.Set("pageToken", r.PageToken)
	}

	return
}

func (c *Client) PlaylistItemsList(ctx context.Context, req PlaylistItemsListRequest) (*response.PlaylistItemsListResponse, error) {

	// send request
	resp, err := c.doGetRequest(ctx, playlistItemsListEndpoint, req.Values())
	if err != nil {
		return nil, err
	}

	// parse body
	body, err := c.parseBody(resp)
	if err != nil {
		return nil, err
	}

	// unmarshal body
	data := response.PlaylistItemsListResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
