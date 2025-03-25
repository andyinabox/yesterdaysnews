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

	// if there are warnings or errors they will be prepended to the json string, so we need to separate them
	// NOTE: this is probably because I'm combining stdout and stderr :(
	msg, data, found := strings.Cut(string(result), "{")
	// if we have content before the first '{', assume it's an error or warning message
	// and set the result to the part after that
	if found && msg != "" {
		result = []byte("{" + data)
		log.Warn("youtubedownloader: " + msg)
	}

	if err != nil {

		// handle requested format not available errors
		if strings.Contains(strings.ToLower(msg), "requested format is not available") {
			err = ErrRequestedFormatNotAvailable
		}

		err = fmt.Errorf("error executing youtubedownloader: %w", err)

	}

	// log.Debug(string(result))

	return result, err

}
