package mediatool

import (
	"context"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/shellargs"
)

func (t *Tool) CutVideo(ctx context.Context, input, output string, start, duration Duration) error {

	options := shellargs.New().
		AddKeyed("-i", input).
		AddKeyed("-ss", start.String()).
		AddKeyed("-to", (start+duration).String()).
		AddKeyed("-c", "copy").
		Add(output)

	_, err := t.executeFfmpeg(ctx, options)

	// slog.Debug("ffmpeg cut video", "result", string(result))

	return err
}
