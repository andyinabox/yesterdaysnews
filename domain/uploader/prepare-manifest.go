package uploader

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Manifest struct {
	Date  time.Time     `json:"date"`
	Files ManifestFiles `json:"files"`
}

type ManifestFiles struct {
	VideoFile    string   `json:"video"`
	SubsFile     string   `json:"subs"`
	CombinedFile string   `json:"combined"`
	ModelFile    string   `json:"model"`
	Clips        []string `json:"clips"`
}

func (u *Uploader) CreateManifest(ctx context.Context, dir string) (*Manifest, error) {
	manifest := &Manifest{
		Date: time.Now(),
	}

	if checkFileExists(filepath.Join(dir, u.cfg.VideoFileName)) {
		manifest.Files.VideoFile = u.cfg.VideoFileName
	} else {
		return nil, fmt.Errorf("video file %q missing", u.cfg.VideoFileName)
	}

	if checkFileExists(filepath.Join(dir, u.cfg.SubsFileName)) {
		manifest.Files.SubsFile = u.cfg.SubsFileName
	} else {
		// return nil, fmt.Errorf("video file %q missing", u.cfg.SubsFileName)
	}

	if checkFileExists(filepath.Join(dir, u.cfg.CombinedFileName)) {
		manifest.Files.CombinedFile = u.cfg.CombinedFileName
	} else {
		// return nil, fmt.Errorf("video file %q missing", u.cfg.CombinedFileName)
	}

	if checkFileExists(filepath.Join(dir, u.cfg.ModelFileName)) {
		manifest.Files.ModelFile = u.cfg.ModelFileName
	} else {
		return nil, fmt.Errorf("video file %q missing", u.cfg.ModelFileName)
	}

	entries, err := os.ReadDir(filepath.Join(dir, u.cfg.ClipsDirName))
	if err != nil {
		return nil, err
	}

	manifest.Files.Clips = make([]string, len(entries))

	for i, f := range entries {
		manifest.Files.Clips[i] = filepath.Join(u.cfg.ClipsDirName, f.Name())
	}

	return manifest, nil
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}
