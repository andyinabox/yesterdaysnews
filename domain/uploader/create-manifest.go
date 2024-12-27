package uploader

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Manifest struct {
	Date  time.Time     `json:"date"`
	ID    string        `json:"id"`
	Files ManifestFiles `json:"files"`
}

type ManifestFiles struct {
	VideoFile string `json:"video"`
	// SubsFile     string   `json:"subs"`
	// CombinedFile string   `json:"combined"`
	ModelFile string   `json:"model"`
	Clips     []string `json:"clips"`
}

func (u *Uploader) CreateManifest(ctx context.Context, dir string) (*Manifest, error) {

	now := time.Now()

	manifest := &Manifest{
		ID:   timestampFromTime(now),
		Date: now,
	}

	if checkFileExists(filepath.Join(dir, u.cfg.VideoFileName)) {
		manifest.Files.VideoFile = u.cfg.VideoFileName
	} else {
		return nil, fmt.Errorf("video file %q missing", u.cfg.VideoFileName)
	}

	// if checkFileExists(filepath.Join(dir, u.cfg.SubsFileName)) {
	// 	manifest.Files.SubsFile = u.cfg.SubsFileName
	// } else {
	// 	return nil, fmt.Errorf("video file %q missing", u.cfg.SubsFileName)
	// }

	// if checkFileExists(filepath.Join(dir, u.cfg.CombinedFileName)) {
	// 	manifest.Files.CombinedFile = u.cfg.CombinedFileName
	// } else {
	// 	return nil, fmt.Errorf("video file %q missing", u.cfg.CombinedFileName)
	// }

	if checkFileExists(filepath.Join(dir, u.cfg.ModelFileName)) {
		manifest.Files.ModelFile = u.cfg.ModelFileName
	} else {
		return nil, fmt.Errorf("video file %q missing", u.cfg.ModelFileName)
	}

	entries, err := filepath.Glob(filepath.Join(dir, u.cfg.ClipsDirName, "*.webm"))
	if err != nil {
		return nil, err
	}

	manifest.Files.Clips = make([]string, len(entries))

	for i, path := range entries {
		manifest.Files.Clips[i] = strings.Replace(path, dir+"/", "", 1)
	}

	return manifest, nil
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}
