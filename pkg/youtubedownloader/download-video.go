package youtubedownloader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// DownloadVideo will execute a download with JSON response, so it's possible to get info about
// the downloaded video
func (c *Client) DownloadVideo(ctx context.Context, url string, req Request) (*YouTubeVideo, error) {

	req.DumpJSON = true
	req.NoSimulate = true

	result, err := c.Execute(ctx, url, req)
	if err != nil {
		return nil, fmt.Errorf("error downloading video %q: %w", url, err)
	}

	video := YouTubeVideo{}
	err = json.NewDecoder(bytes.NewReader(result)).Decode(&video)
	if err != nil {

		// if there is an error it will be printed on the first line, follwed by json
		errStr, _, found := strings.Cut(string(result), "\n")
		if !found {
			errStr = ""
		}

		return nil, fmt.Errorf("error unmarshaling video download result: %w: %s", err, errStr)
	}

	return &video, err
}
