package videoeditor

import (
	"context"
	"strings"
	"time"
)

func (e *Editor) GetVideoLength(ctx context.Context, inputPath string) (time.Duration, error) {

	options := map[string]string{
		"-v":              "error",
		"-show_entries":   "",
		"format=duration": "",
		"-of":             "default=noprint_wrappers=1:nokey=1",
	}

	result, err := e.executeFfprobe(ctx, inputPath, options)

	if err != nil {
		return 0, err
	}

	return time.ParseDuration(strings.TrimSpace(string(result)) + "s")
}
