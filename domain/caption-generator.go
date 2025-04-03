package domain

import "time"

type CaptionGenerator interface {
	Caption(string) string
	Subtitles(d time.Duration) ([]byte, error)
	MinCaptionLength() int
	MaxCaptionLength() int
}
