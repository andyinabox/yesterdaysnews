package mediatool

import (
	"context"

	"code.andydayton.com/andy/yesterdaysnews/pkg/shellargs"
)

func (t *Tool) ExtractPNG(ctx context.Context, input, output string, start Duration) error {
	options := shellargs.New().
		AddKeyed("-i", input).
		AddKeyed("-ss", start.String()).
		AddKeyed("-frames:v", "1").
		Add(output)

	_, err := t.executeFfmpeg(ctx, options)

	// slog.Debug("ffmpeg extract png", "result", string(result))

	return err
}
