package domain

import "context"

type StreamErrType string

const (
	StreamErrTODO               = "StreamErrTODO"
	StreamErrGetVideoID         = "StreamErrGetVideoID"
	StreamErrDownloadVideo      = "StreamErrDownloadVideo"
	StreamErrCutVideo           = "StreamErrCutVideo"
	StreamErrGetVideoEditPoints = "StreamErrGetVideoEditPoints"
	StreamErrFatal              = "StreamErrFatal"
)

type StreamErr interface {
	error
	Type() StreamErrType
}

type ClipStreamer interface {
	ErrorStream(context.Context, func(StreamErr)) chan<- StreamErr
	VideoIDStream(context.Context, chan<- StreamErr, ...string) <-chan string
	VideoDownloadStream(context.Context, chan<- StreamErr, <-chan string) <-chan string
	VideoCutStream(context.Context, chan<- StreamErr, <-chan string) <-chan string
	VideoUploadStream(context.Context, chan<- StreamErr, <-chan string) <-chan string
}
