package mediatool

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/shellargs"
)

const VideoListFileName = ".file_list.txt"

func (t *Tool) CombineVideos(ctx context.Context, files []string, outFile string) (file string, err error) {

	log.Debug("combine videos", "files", files, "outFile", outFile)

	if len(files) == 0 {
		err = errors.New("no input files provided")
		return
	}

	// tmp, err := os.CreateTemp("", "videos.*.txt")
	// if err != nil {
	// 	return "", err
	// }
	// defer os.Remove(tmp.Name())

	// generate file with filenames
	fileList := ""
	for _, f := range files {
		// abs, err := filepath.Abs(f)
		// if err != nil {
		// 	return "", fmt.Errorf("cannot find absolute path for %s, %w", f, err)
		// }
		fileList = fileList + fmt.Sprintf("file '%s'\n", f)
	}

	// log.Debug(fileList)

	// _, err = tmp.Write([]byte(fileList))
	// if err != nil {
	// 	return "", fmt.Errorf("error writing temp file: %w", err)
	// }
	// tmp.Close()

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

	log.Debug(string(result))

	return outFile, nil
}
