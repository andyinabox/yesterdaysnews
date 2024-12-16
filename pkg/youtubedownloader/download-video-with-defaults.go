package youtubedownloader

import (
	"context"
	"path/filepath"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/shellargs"
)

func (d *Downloader) DownloadVideoWithDefaults(ctx context.Context, url string, outputDir string) error {

	options := shellargs.
		New().
		AddKeyedSingleQuoted("--format", "bv[ext=mp4][height<=1280]").
		Add("--write-auto-subs").
		AddKeyed("--sub-format", "vtt").
		AddKeyedSingleQuoted("--output", filepath.Join(outputDir, "%(id)s.%(ext)s")).
		Add(url)

	log.Debug("download video with defaults", "url", url, "outputDir", outputDir, "options", options)

	return d.execute(ctx, options)

}
