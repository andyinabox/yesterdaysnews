package builder

import (
	"context"
	"sync"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func (b *Builder) videoCutStream(ctx context.Context, videoFiles <-chan string) <-chan string {
	clipStream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		close(clipStream)
	}

	cutVideo := func(filePath string) {
		defer wg.Done()

		// log.Infof("done cutting clip from %q: %v", filePath, edit)
		editPoints, err := b.vp.GetVideoEditPoints(ctx, filePath, util.Seconds(b.cfg.MinClipLengthSeconds), util.Seconds(b.cfg.MaxClipLengthSeconds))
		if err != nil {
			b.error(domain.ErrTypeCutVideo, err)
			return
		}

		for clip := range b.vp.CutVideoStream(ctx, b.errs, filePath, editPoints) {
			clipStream <- clip
		}
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

				wg.Add(1)
				go cutVideo(filePath)
			}
		}
	}()

	return clipStream
}
