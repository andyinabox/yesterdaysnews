package mediatool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

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

	// slog.Debug("ffprobe get video info", "result", string(result))

	if err != nil {
		return 0, fmt.Errorf("error getting video info: %w", err)
	}

	data := ffprobeVideLengthResp{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		return 0, fmt.Errorf("error unmarshaling video info: %w", err)
	}

	if data.Format.Duration == "" {
		err = errors.New("duration is empty")
		return 0, fmt.Errorf("error getting video duration: %w", err)
	}

	d, err := time.ParseDuration(strings.TrimSpace(data.Format.Duration) + "s")
	if err != nil {
		return 0, fmt.Errorf("error parsing video duration: %w", err)
	}

	return Duration(d), nil
}
