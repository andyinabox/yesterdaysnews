package youtubedownloader

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/shell"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/shellargs"
)

type Client struct {
	binPath string // path to yt-dlp binary
	shell   *shell.Shell
}

func New(binPath string) *Client {

	sh := shell.New()

	if binPath == "" {
		binPath = sh.MustGetBinaryPath("yt-dlp")
	}

	return &Client{
		binPath: binPath,
		shell:   sh,
	}
}

func (c *Client) DownloadVideo(ctx context.Context, req DownloadVideoRequest) error {

	log.Debug("download video", "id", req.VideoID, "request", req)

	result, err := c.shell.Execute(ctx, fmt.Sprintf("%s%s", c.binPath, req.Args()))

	log.Debug(string(result))

	return err

}

type DownloadVideoRequest struct {
	VideoID       string
	Format        string
	WriteAutoSubs bool
	SubFormat     string
	Output        string
}

func (r *DownloadVideoRequest) Args() *shellargs.Args {

	if r.VideoID == "" {
		log.Error("missing VideoID in download request")
	}

	a := shellargs.New()

	if r.Format != "" {
		a.AddKeyedSingleQuoted("--format", r.Format)
	}

	if r.WriteAutoSubs {
		a.Add("--write-auto-subs")
	}

	if r.SubFormat != "" {
		a.AddKeyed("--sub-format", r.SubFormat)
	}

	if r.Output != "" {
		a.AddKeyedSingleQuoted("--output", r.Output)
	}

	a.Add(r.VideoID)

	return a
}
