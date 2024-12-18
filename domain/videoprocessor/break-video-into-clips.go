package videoprocessor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/mediatool"
	"golang.org/x/exp/rand"
)

func (p *Processor) BreakVideoIntoClips(ctx context.Context, inputPath, outDir string, minLength, maxLength int) error {

	maxShells := runtime.NumCPU()
	totalShells := 0

	// get length of input video
	totalDuration, err := p.mt.GetVideoLength(ctx, inputPath)
	if err != nil {
		return err
	}
	log.Debugf("Duration: %s", totalDuration)

	// split filename into parts
	base, ext := getFnParts(inputPath)

	var i int
	var playhead mediatool.Duration
	var wg sync.WaitGroup

	// chop into pieces
	for {
		start := playhead
		duration := getRandDuration(minLength, maxLength)

		// break out if we're past the total duration
		if start >= totalDuration {
			break
		}

		// if we're close, make the final clip
		if start+duration > totalDuration {
			duration = totalDuration - start
		}

		// generate output filenames
		progressFile := filepath.Join(outDir, fmt.Sprintf("%s-%d.progress.%s", base, i, ext))
		outFile := filepath.Join(outDir, fmt.Sprintf("%s-%d.%s", base, i, ext))

		log.Info("start cutting video segment", "outFile", outFile, "start", start, "duration", duration)

		wg.Add(1)
		totalShells += 1
		go func() {
			defer wg.Done()
			// cut the video and output
			err := p.mt.CutVideo(ctx, inputPath, progressFile, start, duration)
			if err != nil {
				log.Error("error cutting video segment", "error", err, "outFile", outFile)
				return
			}

			// rename to finished filename
			err = os.Rename(progressFile, outFile)
			if err != nil {
				return
			}

			log.Info("finished video segment", "outFile", outFile)
		}()

		if totalShells >= maxShells {
			wg.Wait()
			totalShells = 0
		}

		// move playead and increment index
		playhead = start + duration
		i++
	}

	wg.Wait()

	log.Info("done cutting files", "time")

	return nil
}

func getFnParts(path string) (base, ext string) {
	fnParts := strings.Split(filepath.Base(path), ".")
	base = fnParts[0]
	ext = fnParts[1]
	return
}

func getRandDuration(minLength, maxLength int) mediatool.Duration {
	seconds := rand.Intn(maxLength-minLength) + minLength
	return mediatool.Duration(time.Second * time.Duration(seconds))
}
