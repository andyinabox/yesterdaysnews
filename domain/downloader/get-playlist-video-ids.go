package downloader

import (
	"context"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/youtubeapi"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/youtubeapi/response"
)

// we are requesting a large result so we can iterate and filter by date
const getPlaylistVideosForDateMaxPages = 5

func (d *Downloader) GetPlaylistVideoIDs(ctx context.Context, date time.Time, playlistId, pageToken string) (ids []string, nextPageToken string, err error) {
	var resp *response.PlaylistItemsListResponse

	ids = []string{}

	resp, err = d.ytapi.PlaylistItemsList(ctx, youtubeapi.PlaylistItemsListRequest{
		PlaylistId: playlistId,
		Part:       []string{"snippet"},
		MaxResults: 50, // this is the maximum allowed by the API
		PageToken:  pageToken,
	})
	if err != nil {
		return
	}

	nextPageToken = resp.NextPageToken

	for _, item := range resp.Items {
		if util.IsSameDay(date, item.Snippet.PublishedAt) {
			ids = append(ids, item.Snippet.ResourceID.VideoID)
		}
	}

	return
}

// func (d *Downloader) checkVideo(ctx context.Context, item *response.PlaylistItem, date time.Time, done func(), idsChan chan<- string) {
// 	defer done()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 			videoDate := item.Snippet.PublishedAt
// 			videoID := item.Snippet.ResourceID.VideoID

// 			// first filter by date (todo: better way to do date comparison)
// 			if videoDate.Year() != date.Year() || videoDate.Month() != date.Month() || videoDate.Day() != date.Day() {
// 				log.Debugf("skipping video %q because of date %s", videoID, videoDate)
// 				return
// 			}

// 			// get info using yt-dlp and filter based on that
// 			log.Debug("fetching video info using yt-dlp")
// 			var getVideoErr error
// 			videoInfo, getVideoErr := d.ytdl.GetVideoInfo(ctx, videoID, VideoFormatString)
// 			if getVideoErr != nil {
// 				log.Debug("error getting video info, although this might not be an error and just mistmatched formats", "error", getVideoErr, "videoID", videoID)
// 				return
// 			}

// 			if videoInfo.AspectRatio <= 1 {
// 				log.Debugf("skipping video %q because of aspect ratio %f", videoID, videoInfo.AspectRatio)
// 				return
// 			}

// 			if c, ok := videoInfo.AutomaticCaptions["en"]; !ok || len(c) == 0 {
// 				log.Debugf("skipping video %q because of lack of english subtitles", videoID)
// 				return
// 			}

// 			idsChan <- videoID
// 			return
// 		}
// 	}

// }

// func (d *Downloader) getPlaylistVideoIDs(ctx context.Context, playlistId string, date time.Time, maxResults int) (ids []string, err error) {
// 	ctx, cancel := context.WithCancel(ctx)

// 	idsChan := make(chan string)
// 	pageTokenChan := make(chan string)

// 	currentPage := 0
// 	ids = []string{}

// 	cleanup := func() {
// 		cancel()
// 		// is this ok? was getting errors because of sending to closed channels...
// 		// close(idsChan)
// 		// close(pageTokenChan)
// 	}

// 	// var mu sync.Mutex
// 	// var wg sync.WaitGroup

// 	// fetch initial page of results
// 	go d.fetchPage(
// 		ctx,
// 		playlistId,
// 		date,
// 		"",
// 		idsChan,
// 		pageTokenChan,
// 	)

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			err = ctx.Err()
// 			cleanup()
// 			return
// 		case id := <-idsChan:
// 			ids = append(ids, id)
// 			log.Infof("found video: %s, total found: %d", id, len(ids))
// 			if len(ids) >= maxResults {
// 				cleanup()
// 				return
// 			}
// 		case token := <-pageTokenChan:
// 			currentPage++
// 			if currentPage >= getPlaylistVideosForDateMaxPages {
// 				log.Debug("reached max pages, stopping")
// 				cleanup()
// 				return
// 			}
// 			go d.fetchPage(
// 				ctx,
// 				playlistId,
// 				date,
// 				token,
// 				idsChan,
// 				pageTokenChan,
// 			)
// 		}
// 	}

// }

// func (d *Downloader) fetchPage(ctx context.Context, playlistId string, date time.Time, pageToken string, idsChan chan<- string, pageTokensChan chan<- string) {

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return

// 		default:
// 			ctx, cancel := context.WithCancel(ctx)

// 			var wg sync.WaitGroup

// 			eject := func(err error) {
// 				log.Error(err)
// 				cancel()
// 			}

// 			resp, err := d.ytapi.PlaylistItemsList(ctx, youtubeapi.PlaylistItemsListRequest{
// 				PlaylistId: playlistId,
// 				Part:       []string{"snippet"},
// 				MaxResults: 50,
// 				PageToken:  pageToken,
// 			})
// 			if err != nil {
// 				eject(err)
// 				return
// 			}

// 			if len(resp.Items) == 0 {
// 				eject(errors.New("no playlist items in response"))
// 				return
// 			}

// 			for _, item := range resp.Items {
// 				wg.Add(1)
// 				go d.checkVideo(
// 					ctx,
// 					&item,
// 					date,
// 					func() {
// 						wg.Done()
// 					},
// 					idsChan,
// 				)
// 			}

// 			wg.Wait()

// 			pageTokensChan <- resp.NextPageToken
// 			return
// 		}
// 	}

// }
