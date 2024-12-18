package downloader

import (
	"context"
	"errors"
	"fmt"
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

		// expected files
		videoFile := filepath.Join(req.OutputDir, fmt.Sprintf("%s.mp4", id))
		subsFile := filepath.Join(req.OutputDir, fmt.Sprintf("%s.en.vtt", id))

		// download video in goroutine
		go func() {
			defer wg.Done()

			log.Infof("downloading video %s", videoFile)
			err := d.ytdl.DownloadVideo(ctx, youtubedownloader.DownloadVideoRequest{
				VideoID:       id,
				Format:        "bv[ext=mp4][height<=1280]",
				WriteAutoSubs: true,
				SubFormat:     "vtt",
				Output:        videoFile,
			})
			if err != nil {
				log.Error("error downloading video", "error", err, "id", id)
			}

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
