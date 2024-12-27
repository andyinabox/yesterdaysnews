package uploader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"
)

func (u *Uploader) Upload(ctx context.Context, dir string) (string, error) {

	var wg sync.WaitGroup

	manifestFile, err := os.Open(filepath.Join(dir, "manifest.json"))
	if err != nil {
		return "", fmt.Errorf("error opening manifest file: %w", err)
	}

	defer manifestFile.Close()
	data, err := io.ReadAll(manifestFile)
	if err != nil {
		return "", fmt.Errorf("error reading manifest: %w", err)
	}

	manifest := Manifest{}
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return "", fmt.Errorf("error decoding manifest: %w", err)
	}

	deployDir := manifest.ID

	exists, err := u.osclient.ContainerExists(ctx, u.cfg.ContainerName)
	if err != nil {
		return "", err
	}

	if !exists {
		err = u.osclient.CreatePublicContainer(ctx, u.cfg.ContainerName)
		if err != nil {
			return "", fmt.Errorf("error creating container: %w", err)
		}
	}

	// upload file func
	uploadFile := func(path, contentType string, multipart bool) {
		defer wg.Done()

		key := filepath.Join(deployDir, path)
		filePath := filepath.Join(dir, path)

		log.Debugf("begin uploading file %q as %q", filePath, key)

		file, err := os.Open(filePath)
		if err != nil {
			log.Errorf("error opening file %q: %s", filePath, err)
		}

		_, err = u.osclient.UploadFile(
			ctx,
			u.cfg.ContainerName,
			key,
			file,
			contentType,
			multipart,
		)
		if err != nil {
			log.Errorf("error uploading file %q: %s", key, err)
		}

		log.Debugf("finished uploading %q", key)
	}

	// upload model file
	wg.Add(1)
	go uploadFile(manifest.Files.ModelFile, "application/json", false)

	// upload video file
	wg.Add(1)
	go uploadFile(manifest.Files.VideoFile, "video/mp4", true)

	// upload individual clips
	for _, clipPath := range manifest.Files.Clips {
		wg.Add(1)
		go uploadFile(clipPath, "video/webm", false)
	}

	wg.Wait()

	log.Info("uploading manifest")
	_, err = u.osclient.UploadFile(ctx, u.cfg.ContainerName, filepath.Join(deployDir, "manifest.json"), manifestFile, "application/json", false)
	if err != nil {
		return "", fmt.Errorf("error uploading manifest file: %w", err)
	}

	return deployDir, nil
}
