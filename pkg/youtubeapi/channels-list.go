package youtubeapi

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"code.andydayton.com/andy/yesterdaysnews/pkg/youtubeapi/response"
)

const channelsListEndpoint = "channels"

type ChannelsListRequest struct {
	ForHandle string
	Part      []string
}

func (r *ChannelsListRequest) Values() (q url.Values) {
	q = make(url.Values)
	q.Set("part", strings.Join(r.Part, ","))
	q.Set("forHandle", r.ForHandle)
	return
}

func (c *Client) ChannelsList(ctx context.Context, req ChannelsListRequest) (*response.ChannelListResponse, error) {

	// send request
	resp, err := c.doGetRequest(ctx, channelsListEndpoint, req.Values())
	if err != nil {
		return nil, err
	}

	// parse body
	body, err := c.parseBody(resp)
	if err != nil {
		return nil, err
	}

	// unmarshal body
	data := response.ChannelListResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
