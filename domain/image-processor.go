package domain

import "context"

type ImageProcessor interface {
	AverageImages(ctx context.Context, errs chan<- Error, paths []string, outputFile string) (string, error)
	OverlayImages(ctx context.Context, baseFile, overlayFile, outputFile string) (string, error)
}
