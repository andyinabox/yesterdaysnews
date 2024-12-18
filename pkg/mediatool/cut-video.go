package mediatool

import (
	"context"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/shellargs"
)

func (t *Tool) CutVideo(ctx context.Context, input, output string, start, duration Duration) error {

	options := shellargs.New().
		AddKeyed("-i", input).
		AddKeyed("-ss", start.String()).
		AddKeyed("-to", (start+duration).String()).
		AddKeyed("-c", "copy").
		Add(output)

	result, err := t.executeFfmpeg(ctx, options)

	log.Debug(string(result))

	return err
}
