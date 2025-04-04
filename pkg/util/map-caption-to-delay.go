package util

import (
	"strings"
	"time"
)

func MapCaptionToDelay(cap string, minLength, maxLength int, minDelay, maxDelay float64) time.Duration {
	tokens := strings.Split(cap, " ")
	percent := float64(len(tokens)-minLength) / float64(maxLength-minLength)
	seconds := ((maxDelay - minDelay) * percent) + minDelay
	duration := time.Duration(seconds * float64(time.Second))

	// slog.Debug("caption delay", "duration", duration, "percent", percent, "seconds", seconds)

	return duration
}
