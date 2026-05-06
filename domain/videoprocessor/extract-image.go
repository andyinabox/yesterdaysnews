package videoprocessor

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"code.andydayton.com/andy/yesterdaysnews/domain"
	"code.andydayton.com/andy/yesterdaysnews/domain/errorhandler"
	"code.andydayton.com/andy/yesterdaysnews/pkg/mediatool"
)

func (p *Processor) ExtractImage(ctx context.Context, inFile, outFile string, t time.Duration) (string, error) {

	err := p.mt.ExtractPNG(ctx, inFile, outFile, mediatool.Duration(t))
	if err != nil {
		return "", fmt.Errorf("error extracting PNG: %w", err)
	}

	return outFile, nil
}
func (p *Processor) ExtractImagesStream(ctx context.Context, errs chan<- domain.Error, paths <-chan string, t time.Duration) <-chan string {
	stream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		close(stream)
	}

	extractImage := func(path string) {
		defer wg.Done()
		output := strings.Replace(path, filepath.Ext(path), ".png", 1)
		path, err := p.ExtractImage(ctx, path, output, t)
		if err != nil {
			errs <- errorhandler.Err(domain.ErrTypeExtractImage, fmt.Errorf("error extracting image from %q: %w", path, err))
			return
		}
		stream <- output
	}

	go func() {
		defer cleanup()
		for path := range paths {
			select {
			case <-ctx.Done():
				return
			default:
				wg.Add(1)
				go extractImage(path)
			}
		}
	}()

	return stream
}
