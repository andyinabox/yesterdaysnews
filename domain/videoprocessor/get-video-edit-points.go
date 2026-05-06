package videoprocessor

import (
	"context"
	"fmt"
	"time"

	"code.andydayton.com/andy/yesterdaysnews/domain"
	"code.andydayton.com/andy/yesterdaysnews/pkg/mediatool"
	"golang.org/x/exp/rand"
)

const EditPointsDurationTrim = 2 * time.Second

func (p *Processor) GetVideoEditPoints(ctx context.Context, videoFile string, minClipLength, maxClipLength time.Duration) ([]domain.VideoEdit, error) {
	totalDuration, err := p.mt.GetVideoLength(ctx, videoFile)
	// trim a second off the duration, this helps avoid some errors
	trimmedDuration := totalDuration - mediatool.Duration(EditPointsDurationTrim)
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

		if start >= trimmedDuration {
			break
		}

		if start+duration > trimmedDuration {
			break
			// duration = trimmedDuration - start
		}

		// log.Debugf("new edit: %v, %v", start, duration)
		edits = append(edits, newEdit(start, duration))

		playhead = start + duration
	}

	return edits, nil
}
