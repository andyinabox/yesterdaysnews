package downloader

import (
	"context"
	"path/filepath"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/youtubedownloader"
)

func (d *Downloader) DownloadVideo(ctx context.Context, id, outDir string) (string, error) {
	output := filepath.Join(outDir, "%(id)s.%(ext)s")

	video, err := d.ytdl.DownloadVideo(ctx, id, youtubedownloader.Request{
		Format:        VideoFormatString,
		WriteAutoSubs: true,
		SubFormat:     VideoSubFormat,
		Output:        output,
	})
	if err != nil {
		return "", err
	}

	return video.Filename, nil
}
