package videoeditor

import (
	"context"

	"github.com/charmbracelet/log"
)

func (c *Editor) CutVideo(ctx context.Context, input, output string, start, duration Duration) error {

	options := map[string]string{
		"-i":  input,
		"-ss": start.Timestamp(),
		"-to": (start + duration).Timestamp(),
		// "-c":  "copy",
	}

	result, err := c.executeFfmpeg(ctx, output, options)

	log.Debug(string(result))

	return err
}
