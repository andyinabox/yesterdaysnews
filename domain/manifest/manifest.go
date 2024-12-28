package manifest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	VideoFileName    = "yesterdays-news.mp4"
	SubsVileName     = "yesterdays-news.en.vtt"
	CombinedFileName = "yesterdays-news-cc.mp4"
	ModelFileName    = "yesterdays-news.model.json"
	ClipsDirName     = "clips"
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

func Create(dir string) (*Manifest, error) {

	now := time.Now()

	manifest := &Manifest{
		ID:   strconv.FormatInt(now.Unix(), 10),
		Date: now,
	}

	if checkFileExists(filepath.Join(dir, VideoFileName)) {
		manifest.Files.VideoFile = VideoFileName
	} else {
		return nil, fmt.Errorf("video file %q missing", VideoFileName)
	}

	// if checkFileExists(filepath.Join(dir, DefaultSubsFileName)) {
	// 	manifest.Files.SubsFile = DefaultSubsFileName
	// } else {
	// 	return nil, fmt.Errorf("video file %q missing", DefaultSubsFileName)
	// }

	// if checkFileExists(filepath.Join(dir, CombinedFileName)) {
	// 	manifest.Files.CombinedFile = CombinedFileName
	// } else {
	// 	return nil, fmt.Errorf("video file %q missing", CombinedFileName)
	// }

	if checkFileExists(filepath.Join(dir, ModelFileName)) {
		manifest.Files.ModelFile = ModelFileName
	} else {
		return nil, fmt.Errorf("video file %q missing", ModelFileName)
	}

	entries, err := filepath.Glob(filepath.Join(dir, ClipsDirName, "*.webm"))
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
