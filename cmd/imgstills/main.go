package main

import (
	"context"
	"flag"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/mediatool"
)

var clipsDir string
var verbose bool

func init() {
	flag.StringVar(&clipsDir, "d", "dist/clips", "dir to search for video clips")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	err := godotenv.Load()
	if err != nil {
		log.Warnf("error getting .env: %s", err)
	}
}

func main() {
	ctx := context.Background()

	mt := mediatool.New("", "")

	var mu sync.Mutex
	var wg sync.WaitGroup

	clips, err := filepath.Glob(filepath.Join(clipsDir, "*.webm"))
	if err != nil {
		log.Fatal(err)
	}

	editPoint, err := time.ParseDuration("0s")
	if err != nil {
		log.Fatal(err)
	}

	completedImages := []string{}
	for _, clip := range clips {
		log.Infof("process clip %q", clip)
		wg.Add(1)
		go func() {
			defer wg.Done()
			output := filepath.Join(clipsDir, strings.Replace(filepath.Base(clip), filepath.Ext(clip), ".png", 1))
			log.Debug("extract png", "input", clip, "output", output)
			err := mt.ExtractPNG(ctx, clip, output, mediatool.Duration(editPoint))
			if err != nil {
				log.Errorf("error extracting PNG from %s: %s", clip, err)
				return
			}
			mu.Lock()
			completedImages = append(completedImages, output)
			mu.Unlock()
		}()
	}

	log.Info("processing...")
	wg.Wait()
	log.Infof("finished outputting %d images", len(completedImages))

}
