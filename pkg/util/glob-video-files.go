package util

import (
	"fmt"
	"path/filepath"
)

func GlobVideoFiles(dir, base string) ([]string, error) {
	webm, err := filepath.Glob(filepath.Join(dir, base+".webm"))
	if err != nil {
		return nil, fmt.Errorf("error globbing webm files: %w", err)
	}
	mp4, err := filepath.Glob(filepath.Join(dir, base+".mp4"))
	if err != nil {
		return nil, fmt.Errorf("error globbing mp4 files: %w", err)
	}
	return append(webm, mp4...), nil
}
