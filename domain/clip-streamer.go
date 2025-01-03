package domain

import "context"

type ClipStreamer interface {
	VideoIDStream(context.Context, ...string) <-chan string
	VideoDownloadStream(context.Context, <-chan string) <-chan string
	VideoCutStream(context.Context, <-chan string) <-chan string
	VideoUploadStream(context.Context, <-chan string) <-chan string
}
