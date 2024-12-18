package downloader

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/youtubedownloader"
)

type DownloadRequest struct {
	ChannelUsername string
	Date            time.Time
	MaxResults      int
	OutputDir       string
}

type DownloadResult struct {
	Files []DownloadResultFile `json:"files"`
}

type DownloadResultFile struct {
	Video string `json:"video"`
	Subs  string `json:"subs"`
}

func (d *Downloader) DownloadVideosForChannel(ctx context.Context, req DownloadRequest) (result *DownloadResult, err error) {

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

	result = &DownloadResult{
		Files: make([]DownloadResultFile, 0),
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

			drf := DownloadResultFile{}

			// check to see if expected files were written
			videoExists := checkFileExists(videoFile)
			subsExists := checkFileExists(subsFile)

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

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}
