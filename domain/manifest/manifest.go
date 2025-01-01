package manifest

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func Create(dir string) (*domain.Manifest, error) {

	now := time.Now()

	manifest := &domain.Manifest{
		ID:   strconv.FormatInt(now.Unix(), 10),
		Date: now,
	}

	if util.DoesFileExist(filepath.Join(dir, domain.ManifestVideoFileName)) {
		manifest.Files.VideoFile = domain.ManifestVideoFileName
	} else {
		return nil, fmt.Errorf("video file %q missing", domain.ManifestVideoFileName)
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

	if util.DoesFileExist(filepath.Join(dir, domain.ManifestModelFileName)) {
		manifest.Files.ModelFile = domain.ManifestModelFileName
	} else {
		return nil, fmt.Errorf("video file %q missing", domain.ManifestModelFileName)
	}

	entries, err := filepath.Glob(filepath.Join(dir, domain.ManifestClipsDirName, "*.webm"))
	if err != nil {
		return nil, err
	}

	manifest.Files.Clips = make([]string, len(entries))

	for i, path := range entries {
		manifest.Files.Clips[i] = strings.Replace(path, dir+"/", "", 1)
	}

	return manifest, nil
}
