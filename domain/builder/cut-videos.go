package builder

import (
	"context"
	"log/slog"
	"path/filepath"
	"sync"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func (b *Builder) CutVideos(ctx context.Context, paths <-chan string) <-chan string {
	return b.videoCutStream(ctx, paths)
}

func (b *Builder) videoCutStream(ctx context.Context, videoFiles <-chan string) <-chan string {
	clipStream := make(chan string)

	var wg sync.WaitGroup

	cleanup := func() {
		wg.Wait()
		slog.Debug("closing video clip stream")
		close(clipStream)
	}

	cutVideo := func(filePath string) {
		defer wg.Done()

		editPoints, err := b.vp.GetVideoEditPoints(ctx, filePath, util.Seconds(b.cfg.MinClipLengthSeconds), util.Seconds(b.cfg.MaxClipLengthSeconds))

		if err != nil {
			b.eh.Add(domain.ErrTypeCutVideo, err)
			return
		}

		slog.Info("cutting video into clips", "file", filePath, "count", len(editPoints))
		outDir := filepath.Join(b.cfg.OutputDir, domain.ClipsDirName)
		for clip := range b.vp.CutVideoStream(ctx, b.errs, filePath, outDir, editPoints) {
			slog.Info("finished cutting clip", "file", clip)
			clipStream <- clip
		}
		slog.Debug("exiting cutVideo loop")
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
				slog.Info("start cutting clips", "file", filePath)
				go cutVideo(filePath)
			}
		}
	}()

	return clipStream
}
