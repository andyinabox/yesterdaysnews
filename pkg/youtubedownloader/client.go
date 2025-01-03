package youtubedownloader

import (
	"context"
	"errors"
	"fmt"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/shell"
)

var ErrNoURLProvided = errors.New("no url provided")

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

func (c *Client) Execute(ctx context.Context, url string, req Request) ([]byte, error) {

	if url == "" {
		return nil, ErrNoURLProvided
	}

	log.Debug("execute youtubedownloader", "url", url, "request", req)

	result, err := c.shell.Execute(ctx, fmt.Sprintf("%s%s '%s'", c.binPath, req.String(), url))

	if err != nil {
		err = fmt.Errorf("error executing youtubedownloader: %s: %w", string(result), err)
	}

	log.Debug(string(result))

	return result, err

}
