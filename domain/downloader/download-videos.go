package downloader

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/youtubedownloader"
)

type DownloadFile struct {
	Video string `json:"video"`
	Subs  string `json:"subs"`
}

func (d *Downloader) DownloadVideosForUsername(ctx context.Context, userName string, date time.Time, maxResults int) (files []DownloadFile, err error) {

	files = []DownloadFile{}

	var playlistId string
	playlistId, err = d.GetUploadsPlaylistIdForUser(ctx, userName)
	if err != nil {
		return
	}

	var ids []string
	ids, err = d.GetPlaylistVideosForDate(ctx, playlistId, date, maxResults)
	if err != nil {
		return
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, id := range ids {
		wg.Add(1)

		// expected files
		videoFile := filepath.Join(d.cfg.OutputDir, fmt.Sprintf("%s.mp4", id))
		subsFile := filepath.Join(d.cfg.OutputDir, fmt.Sprintf("%s.en.vtt", id))

		// download video in goroutine
		go func() {
			defer wg.Done()
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

			// check to see if expected files were written
			videoExists := checkFileExists(videoFile)
			subsExists := checkFileExists(subsFile)

			if !videoExists {
				log.Errorf("expected file not found: %s", videoFile)
			}
			if !subsExists {
				log.Errorf("expected file not found: %s", subsFile)
			}

			// don't add to files list if either file is missing
			if !videoExists || !subsExists {
				return
			}

			log.Info("downloaded video", "video", videoFile, "subs", subsFile)

			mu.Lock()
			files = append(files, DownloadFile{
				Video: videoFile,
				Subs:  subsFile,
			})
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
