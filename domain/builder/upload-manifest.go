package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (b *Builder) UploadManifest(ctx context.Context, uploadDir string, manifest *domain.Manifest) error {
	log.Info("saving manifest")

	// TODO: check manifest?

	data, err := json.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("error marshaling manifest: %w", err)
	}

	manifestFilePath := filepath.Join(b.cfg.OutputDir, "manifest.json")
	err = os.WriteFile(manifestFilePath, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing manifest file %q: %w", manifestFilePath, err)
	}

	manifestFileKey := filepath.Join(uploadDir, "manifest.json")
	log.Infof("uploading %q as %q", manifestFilePath, manifestFileKey)
	manifestFileKey, err = b.up.UploadFile(ctx, manifestFilePath, manifestFileKey, "application/json", false)
	if err != nil {
		return fmt.Errorf("error uploading manifest file %q: %w", manifestFilePath, err)
	}

	return nil
}
