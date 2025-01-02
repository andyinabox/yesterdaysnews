package videoprocessor

import (
	"context"
	"errors"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (p *Processor) GetVideoEditPoints(ctx context.Context, videoFile string, minClipLength, maxClipLength time.Duration) ([]domain.VideoEdit, error) {
	return nil, errors.New("not implemented")
}
