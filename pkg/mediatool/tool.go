package mediatool

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
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

	log.Debug("create new video editor", "ffmpegPath", ffmpegPath, "ffprobePath", ffprobePath)
	return &Tool{
		ffmpegPath:  ffmpegPath,
		ffprobePath: ffprobePath,
		shell:       sh,
	}
}

func (t *Tool) execute(ctx context.Context, exePath string, options *shellargs.Args) ([]byte, error) {

	command := fmt.Sprintf("%s%s", exePath, options)

	log.Debug(command)

	return t.shell.Execute(ctx, command)
}

func (t *Tool) executeFfmpeg(ctx context.Context, options *shellargs.Args) ([]byte, error) {
	return t.execute(ctx, t.ffmpegPath, options)
}

func (t *Tool) executeFfprobe(ctx context.Context, options *shellargs.Args) ([]byte, error) {
	return t.execute(ctx, t.ffprobePath, options)
}
