package clipstreamer

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (s *Streamer) VideoCutStream(ctx context.Context, errs chan<- domain.StreamErr, videoFiles <-chan string) <-chan string {
	clipStream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		close(clipStream)
	}

	cutVideo := func(i int, filePath string, edit domain.VideoEdit) {
		defer wg.Done()

		// output filename
		ext := filepath.Ext(filePath)
		outPath := strings.Replace(filePath, ext, fmt.Sprintf("-%d%s", i, ext), 1)

		// do edit
		file, err := s.vp.CutVideo(ctx, filePath, outPath, edit)
		if err != nil {
			errs <- NewErr(domain.StreamErrCutVideo, fmt.Errorf("error cutting video %q: %w", filePath, err))
			return
		}

		clipStream <- file
	}

	go func() {
		defer cleanup()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				filePath, open := <-videoFiles

				// stop and cleanup if the channel is closed
				if !open {
					return
				}

				// get edit point
				editPoints, err := s.vp.GetVideoEditPoints(ctx, filePath, s.cfg.MinClipLength, s.cfg.MaxClipLength)

				// handle error and continue in loop
				if err != nil {
					errs <- NewErr(domain.StreamErrGetVideoEditPoints, fmt.Errorf("error getting video %q edit points: %w", filePath, err))
					continue
				}

				// do  edit for each edit point
				wg.Add(len(editPoints))
				for i, ep := range editPoints {
					go cutVideo(i, filePath, ep)
				}

			}

		}
	}()

	return clipStream
}
