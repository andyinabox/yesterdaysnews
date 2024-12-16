package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"golang.org/x/exp/rand"

	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/videoeditor"
)

var verbose bool
var inputPath, outputDir string
var minLength, maxLength int

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&inputPath, "i", "", "path to input video")
	flag.StringVar(&outputDir, "o", "", "path to output dir")
	flag.IntVar(&minLength, "min", 5, "minimum clip length")
	flag.IntVar(&maxLength, "max", 20, "maximum clip length")
	flag.Parse()

	if inputPath == "" || outputDir == "" {
		log.Fatal("must set input and output paths")
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	executionStartTime := time.Now()

	ctx := context.Background()
	editor := videoeditor.New(os.Getenv("FFMPEG_PATH"), os.Getenv("FFPROBE_PATH"))

	totalDuration, err := editor.GetVideoLength(ctx, inputPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Duration: %s", totalDuration)

	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	parts := strings.Split(filepath.Base(inputPath), ".")
	fileNameBase := parts[0]
	fileNameExt := parts[1]

	// var wg sync.WaitGroup

	var i int
	var playhead videoeditor.Duration
	for {
		start := playhead
		duration := getRandDuration()

		// break out if we're past the total duration
		if start >= videoeditor.Duration(totalDuration) {
			break
		}

		// if we're close, make the final clip
		if start+duration > videoeditor.Duration(totalDuration) {
			duration = videoeditor.Duration(totalDuration) - start
		}

		processingOutFile := filepath.Join(outputDir, fmt.Sprintf("%s-%d.progress.%s", fileNameBase, i, fileNameExt))
		outFile := filepath.Join(outputDir, fmt.Sprintf("%s-%d.%s", fileNameBase, i, fileNameExt))

		log.Info("start cutting video segment", "outFile", outFile, "start", start.Timestamp(), "duration", duration.Timestamp())

		// wg.Add(1)
		// go func() {
		// defer wg.Done()
		err := editor.CutVideo(ctx, inputPath, processingOutFile, start, duration)
		if err != nil {
			log.Error("error cutting video segment", "error", err, "outFile", outFile)
			return
		}

		err = os.Rename(processingOutFile, processingOutFile)
		if err != nil {
			log.Error("error renaming video file", "error", err, "outFile", outFile)
		}

		log.Info("finished video segment", "outFile", outFile)
		// }()

		playhead = start + duration
		i++
	}

	// wg.Wait()

	executionDuration := time.Now().Sub(executionStartTime)

	log.Info("done cutting files", "time", executionDuration.String())

}

func getRandDuration() videoeditor.Duration {
	seconds := rand.Intn(maxLength-minLength) + minLength
	return videoeditor.Duration(time.Second * time.Duration(seconds))
}
