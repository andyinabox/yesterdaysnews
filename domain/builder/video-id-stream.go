package builder

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (b *Builder) videoIDStream(ctx context.Context, date time.Time, playlistIds []string) <-chan string {
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

					log.Infof("fetch video ids for %q", playlistId)

					// get a fresh batch of video ids
					ids, pageToken, err = b.dl.GetPlaylistVideoIDs(ctx, date, playlistId, pageToken)

					if err != nil {
						b.error(domain.ErrTypeGetVideoID, err)
					}

					for _, id := range ids {
						if id == "" {
							err = fmt.Errorf("recieved empty video ID for playlist %q (pageToken %q)", playlistId, pageToken)
							b.error(domain.ErrTypeGetVideoID, err)
							continue
						}

						log.Infof("found valid video id: %q", id)

						idStream <- id
						count++
						if count >= b.cfg.DownloadCountPerPlaylist {
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
