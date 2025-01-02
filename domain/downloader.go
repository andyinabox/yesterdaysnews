package domain

import (
	"context"
	"time"
)

// type DownloadRequest struct {
// 	ChannelUsername string
// 	Date            time.Time
// 	MaxResults      int
// 	OutputDir       string
// }

// type DownloadResult struct {
// 	Files []DownloadResultFile `json:"files"`
// }

// type DownloadResultFile struct {
// 	Video string `json:"video"`
// 	Subs  string `json:"subs"`
// }

type Downloader interface {
	DownloadVideo(ctx context.Context, id, outDir string) (string, error)
	// DownloadVideosForChannel(context.Context, DownloadRequest) (*DownloadResult, error)
	GetPlaylistVideoIDs(ctx context.Context, date time.Time, playlistId, pageToken string) (ids []string, nextPageToken string, err error)
	GetChannelPlaylistID(ctx context.Context, channelName string) (string, error)
}
