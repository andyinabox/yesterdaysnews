package mediatool

import (
	"context"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/shellargs"
)

func (t *Tool) AddCaptions(ctx context.Context, videoFile, subsFile, outFile string) (string, error) {
	options := shellargs.New().
		AddKeyedSingleQuoted("-i", videoFile).
		AddKeyedSingleQuoted("-i", subsFile).
		AddKeyed("-c", "copy").
		AddKeyed("-c:s", "mov_text").
		Add(outFile)

	result, err := t.executeFfmpeg(ctx, options)
	if err != nil {
		return "", err
	}

	log.Debug(string(result))

	return outFile, nil

}
