package domain

import (
	"context"
	"time"
)

// FilterReason identifies at which stage of the playlist video ID filter
// pipeline a candidate was handled. Every item emits FilterReasonFetched plus
// exactly one terminal reason.
type FilterReason string

const (
	FilterReasonFetched           FilterReason = "fetched"
	FilterReasonWrongDate         FilterReason = "wrongDate"
	FilterReasonFormatUnavailable FilterReason = "formatUnavailable"
	FilterReasonTooLarge          FilterReason = "tooLarge"
	FilterReasonNoCaptions        FilterReason = "noCaptions"
	FilterReasonInfoError         FilterReason = "infoError"
	FilterReasonAccepted          FilterReason = "accepted"
)

// FilterCallback is invoked once per playlist item as it moves through the
// filter pipeline.
type FilterCallback func(reason FilterReason)

type YouTubeService interface {
	Setup(ctx context.Context) error
	DownloadVideo(ctx context.Context, id, outDir string) (string, error)
	DownloadVideoStream(ctx context.Context, errs chan<- Error, ids <-chan string, outDir string) <-chan string
	GetPlaylistVideoIDs(ctx context.Context, errs chan<- Error, date time.Time, maxSize uint, playlistId, pageToken string, onFilter FilterCallback) (ids []string, nextPageToken string, err error)
	GetPlaylistVideoIDStream(ctx context.Context, errs chan<- Error, playlistId string, date time.Time, maxSize uint, count int, onFilter FilterCallback) <-chan string
	GetChannelPlaylistID(ctx context.Context, channelName string) (string, error)
}
