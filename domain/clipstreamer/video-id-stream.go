package clipstreamer

import (
	"context"
	"fmt"
	"sync"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (s *Streamer) VideoIDStream(ctx context.Context, errs chan<- domain.StreamErr, playlistIds ...string) <-chan string {
	idStream := make(chan string)

	var wg sync.WaitGroup

	for _, playlistId := range playlistIds {

		// download ids for playlist
		wg.Add(1)
		go func() {
			defer wg.Done()

			var ids []string
			var err error

			// set initial values
			count := 0
			pageToken := ""

			for {
				select {
				case <-ctx.Done():
					return
				default:

					// get a fresh batch of video ids
					ids, pageToken, err = s.dl.GetPlaylistVideoIDs(ctx, s.cfg.VideoDate, playlistId, pageToken)

					if err != nil {
						errs <- NewErr(domain.StreamErrGetVideoID, err)
					}

					for _, id := range ids {
						if id == "" {
							err = fmt.Errorf("recieved empty video ID for playlist %q (pageToken %q)", playlistId, pageToken)
							errs <- NewErr(domain.StreamErrGetVideoID, err)
							continue
						}

						idStream <- id
						count++
						if count >= s.cfg.DownloadCountPerPlaylist {
							return
						}
					}

				}
			}
		}()

	}

	// handle stream closing
	go func() {
		wg.Wait()
		close(idStream)
	}()

	return idStream
}
