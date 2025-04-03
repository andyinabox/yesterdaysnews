package mediatool

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/shellargs"
)

func (t *Tool) Validate(ctx context.Context, path string) error {
	options := shellargs.New().
		AddKeyed("-v", "error").
		Add(path)

	result, err := t.executeFfprobe(ctx, options)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	resultStr := string(result)
	if resultStr != "" {
		return fmt.Errorf("validation error: %s", resultStr)
	}

	return nil
}

func (t *Tool) validateMultiple(ctx context.Context, files []string) (validFiles []string) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	validFiles = []string{}

	wg.Add(len(files))
	for _, fn := range files {
		go func() {
			defer wg.Done()

			err := t.Validate(ctx, fn)
			if err != nil {
				slog.Warn("video file is invalid, skipping", "file", fn, "error", err)
				return
			}

			mu.Lock()
			validFiles = append(validFiles, fn)
			mu.Unlock()
		}()
	}
	wg.Wait()

	return
}
