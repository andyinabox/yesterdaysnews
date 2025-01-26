package builder

import (
	"context"
	"sync"
	"time"

	"github.com/charmbracelet/log"
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
			log.Infof("getting video id stream for %q", playlistID)
			playlistIDs := b.yt.GetPlaylistVideoIDStream(ctx, b.errs, playlistID, date, b.cfg.MaxVideoSize, b.cfg.DownloadCountPerPlaylist)
			for id := range playlistIDs {
				log.Infof("got new video ID: %s", id)
				stream <- id
			}
		}()
	}

	// close stream once all playlistIDs are gathered
	go func() {
		wg.Wait()
		log.Debug("closing video id stream")
		close(stream)
	}()

	return stream
}
