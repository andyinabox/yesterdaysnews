package youtubedownloader

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
)

// would it be better to do this with a single yt-dlp shell and a list of ids?
func (d *Downloader) DownloadVideoListWithDefaults(ctx context.Context, urls []string, outputDir string) error {
	var err error

	errs := []error{}

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := d.DownloadVideoWithDefaults(ctx, url, outputDir)
			if err != nil {
				errs = append(errs, err)
				log.Error(err)
			}
		}()
	}

	wg.Wait()

	if len(errs) != 0 {
		errMsg := ""
		for i, e := range errs {
			errMsg = fmt.Sprintf("%s; %d: %s", errMsg, i, e)
		}
		err = errors.New(errMsg)
	}

	return err
}
