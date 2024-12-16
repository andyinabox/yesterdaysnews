package youtubeapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
)

// we are requesting a large result so we can iterate and filter by date
const getPlaylistVideosForDateResultsCount = 200

type resourceID struct {
	VideoID string `json:"videoId"`
}

type snippet struct {
	PublishedAt time.Time  `json:"publishedAt"`
	ResourceID  resourceID `json:"resourceId"`
}

type getPlaylistVideosForDateRespItem struct {
	Kind    string  `json:"kind"`
	Snippet snippet `json:"snippet"`
}

type getPlaylistVideosForDateResp struct {
	Kind  string                             `json:"kind"`
	Items []getPlaylistVideosForDateRespItem `json:"items"`
}

func (c *Client) GetPlaylistVideosForDate(ctx context.Context, playlistId string, date time.Time, maxResults int) (ids []string, err error) {
	q := make(url.Values)
	q.Add("part", "snippet")
	q.Add("playlistId", playlistId)
	q.Add("maxResults", strconv.Itoa(getPlaylistVideosForDateResultsCount))

	var resp *http.Response
	resp, err = c.doGetRequest(ctx, "playlistItems", q)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Debug(string(body))

	data := getPlaylistVideosForDateResp{}
	err = json.Unmarshal(body, &data)

	if len(data.Items) == 0 {
		err = errors.New("no items in response")
		log.Error(err.Error(), "body", data)
		return
	}

	ids = []string{}

	for _, item := range data.Items {
		d := item.Snippet.PublishedAt
		// log.Debugf("\n%s\n%s\n", date, d)
		if d.Year() == date.Year() && d.Month() == date.Month() && d.Day() == date.Day() {
			// log.Debug("found an video", "id", item.ID)
			ids = append(ids, item.Snippet.ResourceID.VideoID)
		}
		if len(ids) >= maxResults {
			break
		}
	}

	return
}
