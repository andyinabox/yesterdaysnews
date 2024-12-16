package videoeditor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"golang.org/x/exp/rand"
)

func getRandDuration(minLength, maxLength int) Duration {
	seconds := rand.Intn(maxLength-minLength) + minLength
	return Duration(time.Second * time.Duration(seconds))
}

func (e *Editor) BreakVideoIntoClips(ctx context.Context, inputPath, outDir string, minLength, maxLength int) error {

	totalDuration, err := e.GetVideoLength(ctx, inputPath)
	if err != nil {
		return err
	}

	log.Debugf("Duration: %s", totalDuration)

	fnPars := strings.Split(filepath.Base(inputPath), ".")
	fileNameBase := fnPars[0]
	fileNameExt := fnPars[1]

	var i int
	var playhead Duration
	for {
		start := playhead
		duration := getRandDuration(minLength, maxLength)

		// break out if we're past the total duration
		if start >= Duration(totalDuration) {
			break
		}

		// if we're close, make the final clip
		if start+duration > Duration(totalDuration) {
			duration = Duration(totalDuration) - start
		}

		processingOutFile := filepath.Join(outDir, fmt.Sprintf("%s-%d.progress.%s", fileNameBase, i, fileNameExt))
		outFile := filepath.Join(outDir, fmt.Sprintf("%s-%d.%s", fileNameBase, i, fileNameExt))

		log.Info("start cutting video segment", "outFile", outFile, "start", start, "duration", duration)

		err := e.CutVideo(ctx, inputPath, processingOutFile, start, duration)
		if err != nil {
			log.Error("error cutting video segment", "error", err, "outFile", outFile)
			return err
		}

		err = os.Rename(processingOutFile, outFile)
		if err != nil {
			return err
		}

		log.Info("finished video segment", "outFile", outFile)

		playhead = start + duration
		i++
	}

	log.Info("done cutting files", "time")

	return nil
}
