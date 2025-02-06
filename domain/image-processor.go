package domain

import (
	"context"
	"io"
)

type ImageProcessor interface {
	AverageImages(ctx context.Context, errs chan<- Error, paths []string, width, height int, outputFile string) (string, error)
	OverlayImages(ctx context.Context, baseFile, overlayFile io.Reader, outputFile string) (string, error)
}
