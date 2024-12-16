package videoeditor

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/shellargs"
)

type ffprobeFormat struct {
	Duration string `json:"duration"`
}

type ffprobeVideLengthResp struct {
	Format ffprobeFormat `json:"format"`
}

func (e *Editor) GetVideoLength(ctx context.Context, inputPath string) (time.Duration, error) {

	options := shellargs.New().
		Add("-show_format").
		AddKeyed("-of", "json").
		AddKeyed("-v", "error").
		Add(inputPath)

	result, err := e.executeFfprobe(ctx, options)

	log.Debug(string(result))

	if err != nil {
		log.Error("error getting video info", "result", string(result))
		return 0, err
	}

	data := ffprobeVideLengthResp{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		log.Error("error unmarshaling data", "result", string(result))
		return 0, err
	}

	if data.Format.Duration == "" {
		err = errors.New("duration is empty")
		log.Error(err)
		return 0, err
	}

	return time.ParseDuration(strings.TrimSpace(data.Format.Duration) + "s")
}
