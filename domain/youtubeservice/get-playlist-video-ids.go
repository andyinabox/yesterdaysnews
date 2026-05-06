package youtubeservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"code.andydayton.com/andy/yesterdaysnews/domain"
	"code.andydayton.com/andy/yesterdaysnews/domain/errorhandler"
	"code.andydayton.com/andy/yesterdaysnews/pkg/util"
	"code.andydayton.com/andy/yesterdaysnews/pkg/youtubeapi"
	"code.andydayton.com/andy/yesterdaysnews/pkg/youtubeapi/response"
	"code.andydayton.com/andy/yesterdaysnews/pkg/youtubedownloader"
)

func (s *Service) GetPlaylistVideoIDs(ctx context.Context, errs chan<- domain.Error, date time.Time, maxSize uint, playlistId, pageToken string) (ids []string, nextPageToken string, err error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var resp *response.PlaylistItemsListResponse

	ids = []string{}

	// fetch a page of video ids
	resp, err = s.ytapi.PlaylistItemsList(ctx, youtubeapi.PlaylistItemsListRequest{
		PlaylistId: playlistId,
		Part:       []string{"snippet"},
		MaxResults: 50, // this is the maximum allowed by the API
		PageToken:  pageToken,
	})
	if err != nil {
		return
	}
	nextPageToken = resp.NextPageToken

	throttler := time.Tick(s.cfg.ThrottleDownloadsBy)

	slog.Info("recieved playlist items from YouTube API", "count", len(resp.Items), "page", pageToken)

	wg.Add(len(resp.Items))
	for _, item := range resp.Items {
		go func() {
			defer wg.Done()

			id := item.Snippet.ResourceID.VideoID

			// check date
			if !util.IsSameDay(date, item.Snippet.PublishedAt) {
				slog.Debug("skipping video: wrong date", "id", id, "date", item.Snippet.PublishedAt, "targetDate", date)
				return
			}

			if throttler != nil {
				slog.Debug("throttling YouTube ID check", "time", s.cfg.ThrottleDownloadsBy)
				<-throttler
			}

			// this will error if the video format is not available
			videoInfo, err := s.ytdl.GetVideoInfo(ctx, id, VideoFormatString)
			if err != nil {
				if errors.Is(err, youtubedownloader.ErrRequestedFormatNotAvailable) {
					slog.Debug("skipping video because requested format is not available", "id", id)
					return
				}

				errs <- errorhandler.Err(domain.ErrTypeGetVideoID, fmt.Errorf("skipping video %q: error getting video info: %s", id, err))
				return
			}

			// check filesize
			if videoInfo.FilesizeApprox > maxSize {
				slog.Debug("skipping video: size is too large", "id", id, "size", videoInfo.FilesizeApprox)
				return
			}

			// check for captions
			if c, ok := videoInfo.AutomaticCaptions["en"]; !ok || len(c) == 0 {
				slog.Debug("skipping video: no subtitles", "id", id)
				return
			}

			// add to list of ids
			mu.Lock()
			ids = append(ids, id)
			mu.Unlock()
		}()
	}

	wg.Wait()

	return
}

func (s *Service) GetPlaylistVideoIDStream(ctx context.Context, errs chan<- domain.Error, playlistId string, date time.Time, maxSize uint, count int) <-chan string {
	stream := make(chan string)

	cleanup := func() {
		close(stream)
	}

	totalRequests := 0

	go func() {
		defer cleanup()

		var total int
		var err error
		var pageToken string
		var ids []string

		for {
			select {
			case <-ctx.Done():
				return
			default:
				slog.Info("fetch video ids", "playlistId", playlistId)

				if totalRequests >= s.cfg.MaxPlaylistRequests {
					slog.Warn("reached max requests for playlist", "maxRequests", s.cfg.MaxPlaylistRequests, "playlistId", playlistId)
					return
				}

				ids, pageToken, err = s.GetPlaylistVideoIDs(ctx, errs, date, maxSize, playlistId, pageToken)
				totalRequests++
				if err != nil {
					errs <- errorhandler.Err(domain.ErrTypeGetVideoID, err)
					continue
				}

				for _, id := range ids {
					slog.Info("found valid video id", "id", id)
					stream <- id
					total++
					if total >= count {
						return
					}
				}
			}
		}
	}()

	return stream
}

// func (s *Service) checkVideo(ctx context.Context, item *response.PlaylistItem, date time.Time, done func(), idsChan chan<- string) {
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

// func (s *Service) getPlaylistVideoIDs(ctx context.Context, playlistId string, date time.Time, maxResults int) (ids []string, err error) {
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

// func (s *Service) fetchPage(ctx context.Context, playlistId string, date time.Time, pageToken string, idsChan chan<- string, pageTokensChan chan<- string) {

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
