package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/captionschain"
)

func (b *Builder) Model(ctx context.Context, uploadDir string, corpi []domain.Corpus) (string, error) {
	cc := captionschain.New(b.cfg.CaptionPrefixLength)
	cc.BuildFromMultiple(corpi)

	data, err := cc.Save()
	if err != nil {
		return "", fmt.Errorf("error saving model data: %w", err)
	}

	modelFilePath := filepath.Join(b.cfg.OutputDir, domain.ModelFileName)
	log.Infof("saving model file to %q", modelFilePath)
	err = os.WriteFile(modelFilePath, data, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error saving model file %s: %w", domain.ModelFileName, err)
	}

	if b.cfg.SkipUpload {
		log.Info("SkipUpload is true, skipping model upload")
		return strings.TrimPrefix(modelFilePath, b.cfg.OutputDir+"/"), nil
	}

	modelFileKey := filepath.Join(uploadDir, domain.ModelFileName)
	log.Infof("uploading %q as %q", modelFilePath, modelFileKey)
	modelFileKey, err = b.cs.UploadFile(ctx, modelFilePath, modelFileKey, "application/json", false)
	if err != nil {
		return "", fmt.Errorf("error uploading %q: %w", modelFileKey, err)
	}

	return domain.ModelFileName, nil
}
