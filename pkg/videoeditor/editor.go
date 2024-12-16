package videoeditor

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/shell"
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

func (e *Editor) execute(ctx context.Context, exePath string, finalArg string, options map[string]string) ([]byte, error) {
	command := exePath

	// add arguments
	for k, v := range options {
		command = command + " " + k
		if v != "" {
			command = command + " '" + v + "'"
		}
	}

	// add final argument
	if finalArg != "" {
		command = fmt.Sprintf("%s %s", command, finalArg)
	}

	return e.shell.Execute(ctx, command)
}

func (e *Editor) executeFfmpeg(ctx context.Context, outPath string, options map[string]string) ([]byte, error) {
	return e.execute(ctx, e.ffmpegPath, outPath, options)
}

func (e *Editor) executeFfprobe(ctx context.Context, inPath string, options map[string]string) ([]byte, error) {
	return e.execute(ctx, e.ffprobePath, inPath, options)
}
