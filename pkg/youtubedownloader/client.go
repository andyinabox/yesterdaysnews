package youtubedownloader

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/shell"
)

var (
	ErrNoURLProvided               = errors.New("no url provided")
	ErrRequestedFormatNotAvailable = errors.New("requested format is not available")
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

	log.Debug("create new youtubedownloader", "binPath", binPath)

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

	result, err := c.shell.Execute(ctx, fmt.Sprintf("%s%s -- '%s'", c.binPath, req.String(), url))

	if err != nil {

		// handle requested format not available errors
		if strings.Contains(strings.ToLower(string(result)), "requested format is not available") {
			err = ErrRequestedFormatNotAvailable
		}

		// err = fmt.Errorf("error executing youtubedownloader: %s: %w", string(result), err)
		err = fmt.Errorf("error executing youtubedownloader: %w", err)

	}

	log.Debug(string(result))

	return result, err

}
