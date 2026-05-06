package captiongenerator

import (
	"time"

	"code.andydayton.com/andy/yesterdaysnews/pkg/srt"
	"code.andydayton.com/andy/yesterdaysnews/pkg/util"
)

// Subtitles will generate randomized subtitles for the given duration of time
func (g *Generator) Subtitles(targetDuration time.Duration, minCaptionDuration, maxCaptionDuration float64) ([]byte, error) {

	subs := srt.New()

	var previousCaption string
	var duration time.Duration
	var current, end time.Time

	end = end.Add(targetDuration)

	for {
		caption := g.Caption(previousCaption)
		if caption == "" {
			previousCaption = ""
			continue
		}
		duration = util.MapCaptionToDelay(caption, g.cfg.MinCaptionLength, g.cfg.MaxCaptionLength, minCaptionDuration, maxCaptionDuration)

		// slog.Debug("check subtitles length", "current", current.Add(duration), "target", end)
		if current.Add(duration).After(end) {
			break
		}

		// slog.Debug("add caption to subtitles", "duration", duration, "text", caption)
		current = subs.AddToEnd(duration, caption)

		previousCaption = caption
	}

	return []byte(subs.String()), nil
}
