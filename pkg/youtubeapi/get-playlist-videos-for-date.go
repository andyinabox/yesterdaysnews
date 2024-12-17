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
const getPlaylistVideosForDateMaxPages = 5

type thumbnailData struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type resourceID struct {
	VideoID string `json:"videoId"`
}

type pageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type snippet struct {
	PublishedAt time.Time                `json:"publishedAt"`
	ResourceID  resourceID               `json:"resourceId"`
	Thumbnails  map[string]thumbnailData `json:"thumbnails"`
}

type getPlaylistVideosForDateRespItem struct {
	Kind    string  `json:"kind"`
	Snippet snippet `json:"snippet"`
}

type getPlaylistVideosForDateResp struct {
	Kind          string                             `json:"kind"`
	Etag          string                             `json:"etag"`
	NextPageToken string                             `json:"nextPageToken"`
	PrevPageToken string                             `json:"prevPageToken"`
	Items         []getPlaylistVideosForDateRespItem `json:"items"`
	PageInfo      pageInfo                           `json:"pageInfo"`
}

func (c *Client) GetPlaylistVideosForDate(ctx context.Context, playlistId string, date time.Time, maxResults int) (ids []string, err error) {
	q := make(url.Values)
	q.Add("part", "snippet")
	q.Add("playlistId", playlistId)
	q.Add("maxResults", strconv.Itoa(50)) // this is the max allowed

	currentPage := 0
	nextPageToken := ""

	for {

		if currentPage >= getPlaylistVideosForDateMaxPages {
			log.Debug("reached max result pages, finishing")
			return
		}

		if nextPageToken != "" {
			q.Set("pageToken", nextPageToken)
		}

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

	resultsloop:
		for _, item := range data.Items {
			d := item.Snippet.PublishedAt

			// check aspect ratio, filter out vertical videos
			thumb, ok := item.Snippet.Thumbnails["default"]

			if !ok {
				log.Error("no default thumbnail found", "video", item)
				continue resultsloop
			}

			if thumb.Height > thumb.Width {
				log.Debug("portrait video found, skipping")
				continue resultsloop
			}

			// log.Debugf("\n%s\n%s\n", date, d)
			if d.Year() == date.Year() && d.Month() == date.Month() && d.Day() == date.Day() {
				// log.Debug("found an video", "id", item.ID)
				ids = append(ids, item.Snippet.ResourceID.VideoID)
			}

			// finally return results
			if len(ids) >= maxResults {
				return
			}

			nextPageToken = data.NextPageToken
			currentPage++

		}
	}

}
