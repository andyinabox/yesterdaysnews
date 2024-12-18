package youtubedownloader

import (
	"context"
	"encoding/json"
)

func (c *Client) GetVideoInfo(ctx context.Context, url string, format string) (*YouTubeVideo, error) {
	r := Request{
		DumpJSON: true,
		Format:   format,
	}

	result, err := c.Execute(ctx, url, r)
	if err != nil {
		return nil, err
	}

	video := YouTubeVideo{}
	err = json.Unmarshal(result, &video)
	if err != nil {
		return nil, err
	}

	return &video, err
}
