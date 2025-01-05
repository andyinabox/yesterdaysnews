package videoprocessor

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/mediatool"
	"golang.org/x/exp/rand"
)

func (p *Processor) GetVideoEditPoints(ctx context.Context, videoFile string, minClipLength, maxClipLength time.Duration) ([]domain.VideoEdit, error) {
	totalDuration, err := p.mt.GetVideoLength(ctx, videoFile)
	if err != nil {
		return nil, fmt.Errorf("error getting video duration for %q: %w", videoFile, err)
	}

	// log.Debugf("clip length range: %v, %v", minClipLength, maxClipLength)
	// log.Debugf("total duration: %v", totalDuration)

	getRandDuration := func() mediatool.Duration {
		seconds := rand.Intn(int(maxClipLength-minClipLength)) + int(minClipLength)
		return mediatool.Duration(time.Duration(seconds))
	}

	edits := []domain.VideoEdit{}

	var playhead mediatool.Duration

	for {
		start := playhead
		duration := getRandDuration()

		if start >= totalDuration {
			break
		}

		if start+duration > totalDuration {
			duration = totalDuration - start
		}

		// log.Debugf("new edit: %v, %v", start, duration)
		edits = append(edits, newEdit(start, duration))

		playhead = start + duration
	}

	return edits, nil
}
