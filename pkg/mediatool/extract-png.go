package mediatool

import (
	"context"
	"log/slog"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/shellargs"
)

func (t *Tool) ExtractPNG(ctx context.Context, input, output string, start Duration) error {
	options := shellargs.New().
		AddKeyed("-i", input).
		AddKeyed("-ss", start.String()).
		AddKeyed("-frames:v", "1").
		Add(output)

	result, err := t.executeFfmpeg(ctx, options)

	slog.Debug("ffmpeg extract png", "result", string(result))

	return err
}
