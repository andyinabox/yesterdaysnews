package mediatool

import (
	"context"
	"fmt"
	"log/slog"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/shell"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/shellargs"
)

type Tool struct {
	ffmpegPath  string // path to ffmpeg binary
	ffprobePath string // path to ffprobe binary
	shell       *shell.Shell
}

func New(ffmpegPath string, ffprobePath string) *Tool {

	sh := shell.New()

	if ffmpegPath == "" {
		ffmpegPath = sh.MustGetBinaryPath("ffmpeg")
	}

	if ffprobePath == "" {
		ffprobePath = sh.MustGetBinaryPath("ffprobe")
	}

	slog.Debug("create new video editor", "ffmpegPath", ffmpegPath, "ffprobePath", ffprobePath)
	return &Tool{
		ffmpegPath:  ffmpegPath,
		ffprobePath: ffprobePath,
		shell:       sh,
	}
}

func (t *Tool) execute(ctx context.Context, exePath string, options *shellargs.Args) ([]byte, error) {

	command := fmt.Sprintf("%s%s", exePath, options)

	slog.Debug("execute mediatool command", "cmd", command)

	data, err := t.shell.Execute(ctx, command)

	if err != nil {
		err = fmt.Errorf("mediatool error: %w: %s", err, string(data))
	}

	return data, err
}

func (t *Tool) executeFfmpeg(ctx context.Context, options *shellargs.Args) ([]byte, error) {
	return t.execute(ctx, t.ffmpegPath, options)
}

func (t *Tool) executeFfprobe(ctx context.Context, options *shellargs.Args) ([]byte, error) {
	return t.execute(ctx, t.ffprobePath, options)
}
