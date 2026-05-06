package builder

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"code.andydayton.com/andy/yesterdaysnews/domain"
	"code.andydayton.com/andy/yesterdaysnews/domain/captionschain"
)

func (b *Builder) Model(ctx context.Context, uploadDir string, corpi []domain.Corpus) (string, error) {
	cc := captionschain.New(b.cfg.CaptionPrefixLength)

	cc.BuildFromMultiple(corpi)

	data, err := cc.Save()
	if err != nil {
		return "", fmt.Errorf("error saving model data: %w", err)
	}

	modelFilePath := filepath.Join(b.cfg.OutputDir, domain.ModelFileName)
	slog.Info("saving model file", "path", modelFilePath)
	err = os.WriteFile(modelFilePath, data, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error saving model file %s: %w", domain.ModelFileName, err)
	}

	if b.cfg.SkipUpload {
		slog.Info("SkipUpload is true, skipping model upload")
		return strings.TrimPrefix(modelFilePath, b.cfg.OutputDir+"/"), nil
	}

	modelFileKey := filepath.Join(uploadDir, domain.ModelFileName)
	slog.Info("uploading model file", "from", modelFilePath, "to", modelFileKey)
	modelFileKey, err = b.cs.UploadFile(ctx, modelFilePath, modelFileKey, "application/json", false)
	if err != nil {
		return "", fmt.Errorf("error uploading %q: %w", modelFileKey, err)
	}

	return domain.ModelFileName, nil
}
