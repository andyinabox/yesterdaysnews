package videoprocessor

import (
	"context"
	"time"
)

func (p *Processor) GetVideoLength(ctx context.Context, videoFile string) (time.Duration, error) {
	d, err := p.mt.GetVideoLength(ctx, videoFile)
	if err != nil {
		return 0, err
	}

	return time.Duration(d), nil
}
