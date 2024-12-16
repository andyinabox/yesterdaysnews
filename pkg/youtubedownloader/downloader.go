package youtubedownloader

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"

	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/shell"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/shellargs"
)

type Downloader struct {
	path  string // path to yt-dlp binary
	shell *shell.Shell
}

func New(ytdlpPath string) *Downloader {
	return &Downloader{
		path:  ytdlpPath,
		shell: shell.New(),
	}
}

func (d *Downloader) execute(ctx context.Context, options *shellargs.Args) error {
	command := d.path

	command = fmt.Sprintf("%s%s", d.path, options)

	log.Debug("execute command", "comman", command)

	result, err := d.shell.Execute(ctx, command)

	log.Debug(string(result))

	return err
}
