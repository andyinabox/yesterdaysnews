package youtubedownloader

import (
	"context"
	"encoding/json"
)

// DownloadVideo will execute a download with JSON response, so it's possible to get info about
// the downloaded video
func (c *Client) DownloadVideo(ctx context.Context, url string, req Request) (*YouTubeVideo, error) {

	req.DumpJSON = true
	req.NoSimulate = true

	result, err := c.Execute(ctx, url, req)
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
