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

	// var mu sync.Mutex

	for {

		if currentPage >= getPlaylistVideosForDateMaxPages {
			log.Debug("reached max result pages, finishing")
			return
		}

		var resp *response.PlaylistItemsListResponse
		resp, err = d.ytapi.PlaylistItemsList(ctx, youtubeapi.PlaylistItemsListRequest{
			PlaylistId: playlistId,
			Part:       []string{"snippet"},
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
			videoDate := item.Snippet.PublishedAt
			videoID := item.Snippet.ResourceID.VideoID

			// // check aspect ratio, filter out vertical videos
			// thumb, ok := item.Snippet.Thumbnails["standard"]

			// if !ok {
			// 	log.Error("no standard thumbnail found", "video", item)
			// 	continue resultsloop
			// }

			// log.Debug("checking video size", "width", thumb.Width, "height", thumb.Height)
			// if thumb.Height > thumb.Width {
			// 	log.Debug("portrait video found, skipping")
			// 	continue resultsloop
			// }

			// first filter by date (todo: better way to do date comparison)
			if videoDate.Year() != date.Year() || videoDate.Month() != date.Month() || videoDate.Day() != date.Day() {
				log.Debugf("skipping video %q because of date %s", videoID, videoDate)
				continue resultsloop
			}

			// get info using yt-dlp and filter based on that
			log.Debug("fetching video info using yt-dlp")
			var getVideoErr error
			videoInfo, getVideoErr := d.ytdl.GetVideoInfo(ctx, videoID, VideoFormatString)
			if getVideoErr != nil {
				log.Error("error getting video info", "error", getVideoErr, "videoID", videoID)
				continue resultsloop
			}

			if videoInfo.AspectRatio <= 1 {
				log.Debugf("skipping video %q because of aspect ratio %f", videoID, videoInfo.AspectRatio)
				continue resultsloop
			}

			if c, ok := videoInfo.AutomaticCaptions["en"]; !ok || len(c) == 0 {
				log.Debugf("skipping video %q because of lack of english subtitles", videoID)
				continue resultsloop
			}

			// mu.Lock()
			ids = append(ids, videoID)
			log.Debugf("Adding video %s, found %d videos", videoID, len(ids))
			// mu.Unlock()

			// finally return results
			if len(ids) >= maxResults {
				return
			}

			currentPage++

		}
	}

}
