package downloader

import (
	"context"
	"errors"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/youtubeapi"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/youtubeapi/response"
)

// we are requesting a large result so we can iterate and filter by date
const getPlaylistVideosForDateMaxPages = 5

func (d *Downloader) GetPlaylistVideosForDate(ctx context.Context, playlistId string, date time.Time, maxResults int) (ids []string, err error) {
	currentPage := 0
	nextPageToken := ""

	for {

		if currentPage >= getPlaylistVideosForDateMaxPages {
			log.Debug("reached max result pages, finishing")
			return
		}

		var resp *response.PlaylistItemsListResponse
		resp, err = d.ytapi.PlaylistItemsList(ctx, youtubeapi.PlaylistItemsListRequest{
			PlaylistId: playlistId,
			Part:       []string{"snippet", " contentDetails", "status"},
			MaxResults: 50,
			PageToken:  nextPageToken,
		})
		if err != nil {
			return
		}

		nextPageToken = resp.NextPageToken

		if len(resp.Items) == 0 {
			err = errors.New("no playlist items in response")
			return
		}

		ids = []string{}

	resultsloop:
		for _, item := range resp.Items {
			d := item.Snippet.PublishedAt

			// check aspect ratio, filter out vertical videos
			thumb, ok := item.Snippet.Thumbnails["standard"]

			if !ok {
				log.Error("no standard thumbnail found", "video", item)
				continue resultsloop
			}

			log.Debug("checking video size", "width", thumb.Width, "height", thumb.Height)
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

			currentPage++

		}
	}

}

// import (
// 	"context"
// 	"errors"
// 	"net/url"
// 	"strconv"
// 	"time"

// 	"github.com/charmbracelet/log"
// 	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/youtubeapi/response"
// )

// // we are requesting a large result so we can iterate and filter by date
// const getPlaylistVideosForDateMaxPages = 5

// func (c *Client) GetPlaylistVideosForDate(ctx context.Context, playlistId string, date time.Time, maxResults int) (ids []string, err error) {
// 	q := make(url.Values)
// 	q.Add("part", "snippet,contentDetails,status")
// 	q.Add("playlistId", playlistId)
// 	q.Add("maxResults", strconv.Itoa(50)) // this is the max allowed

// 	currentPage := 0
// 	nextPageToken := ""

// 	for {

// 		if currentPage >= getPlaylistVideosForDateMaxPages {
// 			log.Debug("reached max result pages, finishing")
// 			return
// 		}

// 		if nextPageToken != "" {
// 			q.Set("pageToken", nextPageToken)
// 		}

// 		var resp *response.Success
// 		resp, err = c.doGetRequest(ctx, "playlistItems", q)
// 		if err != nil {
// 			return
// 		}

// 		var playlistItems []response.PlaylistItem
// 		playlistItems, err = resp.PlaylistItems()
// 		if err != nil {
// 			return
// 		}

// 		if len(playlistItems) == 0 {
// 			err = errors.New("no playlist items in response")
// 			log.Error(err.Error(), "resp", resp)
// 			return
// 		}

// 		ids = []string{}

// 	resultsloop:
// 		for _, item := range playlistItems {
// 			d := item.Snippet.PublishedAt

// 			// check aspect ratio, filter out vertical videos
// 			thumb, ok := item.Snippet.Thumbnails["standard"]

// 			if !ok {
// 				log.Error("no standard thumbnail found", "video", item)
// 				continue resultsloop
// 			}

// 			log.Debug("checking video size", "width", thumb.Width, "height", thumb.Height)
// 			if thumb.Height > thumb.Width {
// 				log.Debug("portrait video found, skipping")
// 				continue resultsloop
// 			}

// 			// log.Debugf("\n%s\n%s\n", date, d)
// 			if d.Year() == date.Year() && d.Month() == date.Month() && d.Day() == date.Day() {
// 				// log.Debug("found an video", "id", item.ID)
// 				ids = append(ids, item.Snippet.ResourceID.VideoID)
// 			}

// 			// finally return results
// 			if len(ids) >= maxResults {
// 				return
// 			}

// 			nextPageToken = resp.NextPageToken
// 			currentPage++

// 		}
// 	}

// }

// // func getVideoIdsForPlaylist(ctx context.Context, playlistId, pageToken string) (ids []string, nextPageToken string, err error) {
// // 	q := make(url.Values)
// // 	q.Set("part", "snippet,contentDetails,status")
// // 	q.Set("playlistId", playlistId)
// // 	q.Set("maxResults", strconv.Itoa(50)) // this is the max allowed
// // 	if pageToken != "" {
// // 		q.Set("pageToken", pageToken)
// // 	}

// // }
