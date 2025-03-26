package builder

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

func (b *Builder) DownloadVideos(ctx context.Context, date time.Time) <-chan string {
	// get stream of ids
	ids := b.videoIDStream(ctx, date, b.cfg.PlaylistIDs)

	// get stream of downloaded video paths
	return b.yt.DownloadVideoStream(ctx, b.errs, ids, b.cfg.OutputDir)
}

func (b *Builder) videoIDStream(ctx context.Context, date time.Time, playlistIDs []string) <-chan string {

	stream := make(chan string)

	var wg sync.WaitGroup

	// combine multiple playlist streams into single stream
	wg.Add(len(playlistIDs))
	for _, playlistID := range playlistIDs {
		go func() {
			defer wg.Done()
			slog.Info("getting video id stream", "playlistID", playlistID)
			playlistIDs := b.yt.GetPlaylistVideoIDStream(ctx, b.errs, playlistID, date, b.cfg.MaxVideoSize, b.cfg.DownloadCountPerPlaylist)
			for id := range playlistIDs {
				slog.Info("got new video ID", "id", id)
				stream <- id
			}
		}()
	}

	// close stream once all playlistIDs are gathered
	go func() {
		wg.Wait()
		slog.Debug("closing video id stream")
		close(stream)
	}()

	return stream
}
