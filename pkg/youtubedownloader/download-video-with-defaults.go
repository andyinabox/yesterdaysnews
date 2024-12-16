package youtubedownloader

import (
	"context"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func (d *Downloader) DownloadVideoWithDefaults(ctx context.Context, url string, outputDir string) error {

	options := map[string]string{
		"--format":          "bv[ext=mp4][height<=1280]",
		"--write-auto-subs": "",
		"--sub-format":      "vtt",
		"--output":          filepath.Join(outputDir, "%(id)s.%(ext)s"),
	}

	log.Debug("download video with defaults", "url", url, "outputDir", outputDir, "options", options)

	return d.execute(ctx, url, options)

}
