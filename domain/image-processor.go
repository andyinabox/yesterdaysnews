package domain

import (
	"context"
	"io"
)

type ImageProcessor interface {
	AverageImagesStream(ctx context.Context, errs chan<- Error, paths <-chan string, width, height int, outputFile string) (string, error)
	OverlayImages(ctx context.Context, baseFile, overlayFile io.Reader, outputFile string) (string, error)
}
