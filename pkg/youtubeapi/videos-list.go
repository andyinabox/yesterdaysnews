package youtubeapi

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/youtubeapi/response"
)

const videosListEndpoint = "videos"

type VideosListRequest struct {
	Part []string
}

func (r *VideosListRequest) Values() (q url.Values) {
	q = make(url.Values)
	q.Set("part", strings.Join(r.Part, ","))
	return
}

func (c *Client) VideosList(ctx context.Context, req VideosListRequest) (*response.VideosListResponse, error) {

	// send request
	resp, err := c.doGetRequest(ctx, videosListEndpoint, req.Values())
	if err != nil {
		return nil, err
	}

	// parse body
	body, err := c.parseBody(resp)
	if err != nil {
		return nil, err
	}

	// unmarshal body
	data := response.VideosListResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
