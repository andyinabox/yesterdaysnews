package uploader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"
)

func (u *Uploader) Upload(ctx context.Context, dir string) (containerName string, err error) {

	var wg sync.WaitGroup

	manifest, err := u.CreateManifest(ctx, dir)
	if err != nil {
		return "", fmt.Errorf("error creating manifest: %w", err)
	}

	bucketName := u.cfg.BucketNameBase + "-" + manifest.ID
	log.Debugf("creating bucket %q", bucketName)
	err = u.osclient.CreateBucket(ctx, bucketName)
	if err != nil {
		return "", fmt.Errorf("error creating bucket %q: %w", bucketName, err)
	}

	// upload file func
	uploadFile := func(ctx context.Context, bucketName, fileKey, filePath, contentType string, multipart bool) {
		defer wg.Done()

		log.Debugf("begin uploading file %q as %q", filePath, fileKey)

		file, err := os.Open(filePath)
		if err != nil {
			log.Errorf("error opening file %q: %s", filePath, err)
		}

		if multipart {
			_, err = u.osclient.UploadFileMultipart(ctx, bucketName, fileKey, file, contentType)
		} else {
			_, err = u.osclient.UploadFile(ctx, bucketName, fileKey, file, contentType)
		}
		if err != nil {
			log.Errorf("error uploading file %q: %s", fileKey, err)
		}

		log.Debugf("finished uploading %q", fileKey)
	}

	// limiting the number of concurrent uploads
	// maxConcurrent := 20 // runtime.NumCPU()
	// totalConcurrent := 0

	// upload model file
	wg.Add(1)
	// totalConcurrent++
	go uploadFile(ctx, bucketName, manifest.Files.ModelFile, filepath.Join(dir, manifest.Files.ModelFile), "application/json", false)

	// upload video file
	wg.Add(1)
	// totalConcurrent++
	go uploadFile(ctx, bucketName, manifest.Files.VideoFile, filepath.Join(dir, manifest.Files.VideoFile), "video/mp4", true)

	// upload individual clips
	for _, clipPath := range manifest.Files.Clips {
		wg.Add(1)
		// totalConcurrent++
		go uploadFile(ctx, bucketName, clipPath, filepath.Join(dir, clipPath), "video/webm", false)
		// if totalConcurrent >= maxConcurrent {
		// 	wg.Wait()
		// 	totalConcurrent = 0
		// }
	}

	wg.Wait()

	data, err := json.Marshal(*manifest)
	if err != nil {
		return "", fmt.Errorf("error marshaling manifest data: %w", err)
	}

	log.Info("uploading manifest")
	_, err = u.osclient.UploadFile(ctx, bucketName, "manifest.json", bytes.NewReader(data), "application/json")
	if err != nil {
		return "", fmt.Errorf("error uploading manifest file: %w", err)
	}

	return bucketName, nil
}
