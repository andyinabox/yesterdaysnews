package mediatool

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/shellargs"
)

const VideoListFileName = ".file-list.txt"

func (t *Tool) CombineVideos(ctx context.Context, files []string, outFile string) (file string, err error) {

	slog.Debug("combine videos", "count", len(files), "outFile", outFile)

	if len(files) == 0 {
		err = errors.New("no input files provided")
		return
	}

	// generate file with filenames
	var fileList string
	for _, f := range files {
		// ffmpeg format is "file 'path/to/file.mp4'"
		fileList = fileList + fmt.Sprintf("file '%s'\n", f)
	}

	err = os.WriteFile(VideoListFileName, []byte(fileList), os.ModePerm)
	defer os.Remove(VideoListFileName)

	options := shellargs.New().
		AddKeyed("-f", "concat").
		AddKeyedSingleQuoted("-i", VideoListFileName).
		AddKeyed("-c", "copy").
		Add(outFile)

	result, err := t.executeFfmpeg(ctx, options)
	if err != nil {
		return
	}

	slog.Debug("ffmpeg combine videos", "result", string(result))

	return outFile, nil
}
