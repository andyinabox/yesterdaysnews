package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (b *Builder) Manifest(ctx context.Context, uploadDir string, manifest *domain.Manifest) (string, error) {
	log.Info("saving manifest")

	// TODO: check manifest?

	data, err := json.Marshal(manifest)
	if err != nil {
		return "", fmt.Errorf("error marshaling manifest: %w", err)
	}

	manifestFilePath := filepath.Join(b.cfg.OutputDir, domain.ManifestFileName)
	err = os.WriteFile(manifestFilePath, data, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error writing manifest file %q: %w", manifestFilePath, err)
	}

	// finish here if skipping upload step
	if b.cfg.SkipUpload {
		log.Info("SkipUpload is true, skipping manifest upload")
		return strings.TrimPrefix(manifestFilePath, b.cfg.OutputDir+"/"), nil
	}

	manifestFileKey := filepath.Join(uploadDir, domain.ManifestFileName)
	log.Infof("uploading %q as %q", manifestFilePath, manifestFileKey)
	manifestFileKey, err = b.cs.UploadFile(ctx, manifestFilePath, manifestFileKey, "application/json", false)
	if err != nil {
		return "", fmt.Errorf("error uploading manifest file %q: %w", manifestFilePath, err)
	}

	return manifestFileKey, nil
}
