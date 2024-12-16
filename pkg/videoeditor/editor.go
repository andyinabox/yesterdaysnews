package videoeditor

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/shell"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/shellargs"
)

type Editor struct {
	ffmpegPath  string // path to ffmpeg binary
	ffprobePath string // path to ffprobe binary
	shell       *shell.Shell
}

func New(ffmpegPath string, ffprobePath string) *Editor {

	log.Debug("create new video editor", "ffmpegPath", ffmpegPath, "ffprobePath", ffprobePath)

	return &Editor{
		ffmpegPath:  ffmpegPath,
		ffprobePath: ffprobePath,
		shell:       shell.New(),
	}
}

func (e *Editor) execute(ctx context.Context, exePath string, options *shellargs.Args) ([]byte, error) {

	command := fmt.Sprintf("%s%s", exePath, options)

	log.Debug(command)

	return e.shell.Execute(ctx, command)
}

func (e *Editor) executeFfmpeg(ctx context.Context, options *shellargs.Args) ([]byte, error) {
	return e.execute(ctx, e.ffmpegPath, options)
}

func (e *Editor) executeFfprobe(ctx context.Context, options *shellargs.Args) ([]byte, error) {
	return e.execute(ctx, e.ffprobePath, options)
}
