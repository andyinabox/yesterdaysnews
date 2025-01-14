package domain

import (
	"context"
	"time"
)

type YouTubeService interface {
	DownloadVideo(ctx context.Context, id, outDir string) (string, error)
	DownloadVideoStream(ctx context.Context, errs chan<- Error, ids <-chan string, outDir string) <-chan string
	GetPlaylistVideoIDs(ctx context.Context, date time.Time, playlistId, pageToken string) (ids []string, nextPageToken string, err error)
	GetPlaylistVideoIDStream(ctx context.Context, errs chan<- Error, playlistId string, date time.Time, count int) <-chan string
	GetChannelPlaylistID(ctx context.Context, channelName string) (string, error)
}
