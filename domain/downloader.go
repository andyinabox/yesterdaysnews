package domain

import (
	"context"
	"time"
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

type Downloader interface {
	// DownloadVideo() // TODO
	DownloadVideosForChannel(context.Context, DownloadRequest) (*DownloadResult, error)
	GetPlaylistVideosForDate(ctx context.Context, playlistId string, date time.Time, maxResults int) ([]string, error)
	GetUploadsPlaylistIdForChannel(context.Context, string) (string, error)
}
