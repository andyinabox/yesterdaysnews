package videoprocessor

import (
	"context"
	"errors"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (p *Processor) CutVideo(ctx context.Context, inFile, outFile string, edit domain.VideoEdit) (string, error) {
	return "", errors.New("not implemented")
}
