package downloader

import (
	"context"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/youtubedownloader"
)

func (d *Downloader) DownloadVideosForChannel(ctx context.Context, req domain.DownloadRequest) (result *domain.DownloadResult, err error) {

	log.Infof("getting playlistId for channel %q", req.ChannelUsername)
	var playlistId string
	playlistId, err = d.GetUploadsPlaylistIdForChannel(ctx, req.ChannelUsername)
	if err != nil {
		return
	}
	log.Infof("got playlistId %q", playlistId)

	log.Infof("getting ~%d video ids from %s", req.MaxResults, req.Date)
	var ids []string
	ids, err = d.GetPlaylistVideosForDate(ctx, playlistId, req.Date, req.MaxResults)
	if err != nil {
		return
	}
	log.Infof("found %d ids", len(ids))
	log.Debug(strings.Join(ids, "\n"))

	var mu sync.Mutex
	var wg sync.WaitGroup

	result = &domain.DownloadResult{
		Files: make([]domain.DownloadResultFile, 0),
	}

	for _, id := range ids {
		wg.Add(1)

		// download video in goroutine
		go func() {
			defer wg.Done()

			log.Infof("downloading video %s", id)
			video, err := d.ytdl.DownloadVideo(ctx, id, youtubedownloader.Request{
				Format:        VideoFormatString,
				WriteAutoSubs: true,
				SubFormat:     VideoSubFormat,
				Output:        filepath.Join(req.OutputDir, "%(id)s.%(ext)s"),
			})
			if err != nil {
				log.Error("error downloading video", "error", err, "id", id)
			}

			videoFile := video.Filename
			subsFile := strings.Replace(videoFile, filepath.Ext(videoFile), ".en."+VideoSubFormat, 1)

			drf := domain.DownloadResultFile{}

			// check to see if expected files were written
			videoExists := util.DoesFileExist(videoFile)
			subsExists := util.DoesFileExist(subsFile)

			if videoExists {
				drf.Video = videoFile
			} else {
				log.Errorf("expected file not found: %s", videoFile)
			}
			if subsExists {
				drf.Subs = subsFile
			} else {
				log.Errorf("expected file not found: %s", subsFile)
			}

			log.Infof("finished downloading video %s", videoFile)

			mu.Lock()
			result.Files = append(result.Files, drf)
			mu.Unlock()
		}()

	}

	wg.Wait()

	return
}
