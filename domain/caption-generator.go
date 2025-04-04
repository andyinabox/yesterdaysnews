package domain

import "time"

type CaptionGenerator interface {
	Caption(string) string
	Subtitles(length time.Duration, minDuration, maxDuration float64) ([]byte, error)
	MinCaptionLength() int
	MaxCaptionLength() int
}
