package builder

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"code.andydayton.com/andy/yesterdaysnews/domain"
)

func (b *Builder) DownloadVideos(ctx context.Context, date time.Time) <-chan string {
	// get stream of ids
	ids := b.videoIDStream(ctx, date, b.cfg.Playlists)

	// get stream of downloaded video paths
	return b.yt.DownloadVideoStream(ctx, b.errs, ids, b.cfg.OutputDir)
}

func (b *Builder) videoIDStream(ctx context.Context, date time.Time, playlists []domain.PlaylistSource) <-chan string {

	stream := make(chan string)

	var wg sync.WaitGroup

	// combine multiple playlist streams into single stream
	wg.Add(len(playlists))
	for _, playlist := range playlists {
		go func() {
			defer wg.Done()
			slog.Info("getting video id stream", "playlistID", playlist.ID, "name", playlist.Name)
			playlistID := playlist.ID
			onFilter := func(reason domain.FilterReason) {
				b.recordFilterReason(playlistID, reason)
			}
			ids := b.yt.GetPlaylistVideoIDStream(ctx, b.errs, playlist.ID, date, b.cfg.MaxVideoSize, b.cfg.DownloadCountPerPlaylist, onFilter)
			for id := range ids {
				slog.Info("got new video ID", "id", id, "playlistID", playlist.ID)
				b.recordVideoSource(id, playlist.ID)
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
