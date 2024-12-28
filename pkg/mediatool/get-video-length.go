package mediatool

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/shellargs"
)

type ffprobeFormat struct {
	Duration string `json:"duration"`
}

type ffprobeVideLengthResp struct {
	Format ffprobeFormat `json:"format"`
}

func (t *Tool) GetVideoLength(ctx context.Context, inputPath string) (Duration, error) {

	options := shellargs.New().
		Add("-show_format").
		AddKeyed("-of", "json").
		AddKeyed("-v", "error").
		Add(inputPath)

	result, err := t.executeFfprobe(ctx, options)

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

	d, err := time.ParseDuration(strings.TrimSpace(data.Format.Duration) + "s")
	if err != nil {
		return 0, err
	}

	return Duration(d), nil
}
