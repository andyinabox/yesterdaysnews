package videoprocessor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/mediatool"
)

func (p *Processor) CutVideo(ctx context.Context, inFile, outFile string, edit domain.VideoEdit) (string, error) {

	ext := filepath.Ext(outFile)
	progressFile := strings.Replace(outFile, ext, fmt.Sprintf("-progress%s", ext), 1)

	err := p.mt.CutVideo(ctx, inFile, progressFile, mediatool.Duration(edit.Start()), mediatool.Duration(edit.Duration()))
	if err != nil {
		return "", fmt.Errorf("error cutting video segment for %q: %v: %w", inFile, edit, err)
	}

	err = os.Rename(progressFile, outFile)
	if err != nil {
		return "", fmt.Errorf("error renaming %q to %q: %w", progressFile, outFile, err)
	}

	return outFile, nil
}

func (p *Processor) CutVideoStream(ctx context.Context, errs chan<- domain.Error, inFile string, edits []domain.VideoEdit) <-chan string {
	stream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		close(stream)
	}

	cutVideo := func(i int, edit domain.VideoEdit) {
		defer wg.Done()

		ext := filepath.Ext(inFile)
		clipName := strings.Replace(inFile, ext, fmt.Sprintf("-%d%s", i, ext), 1)

		videoClip, err := p.CutVideo(ctx, inFile, clipName, edit)
		if err != nil {
			errs <- errorhandler.Err(domain.ErrTypeCutVideo, err)
			return
		}

		stream <- videoClip
	}

	go func() {
		defer cleanup()
		wg.Add(len(edits))
		for i, edit := range edits {
			go cutVideo(i, edit)
		}
	}()

	return stream
}
