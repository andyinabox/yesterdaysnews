package builder

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (b *Builder) PosterImage(ctx context.Context, uploadDir string, paths []string) (string, error) {

	if b.cfg.OverlayImage == nil {
		return "", errors.New("no overlay image set")
	}

	avgImagePath, err := b.ip.AverageImages(
		ctx,
		b.errs,
		paths,
		domain.VideoWidth,
		domain.VideoHeight,
		filepath.Join(b.cfg.OutputDir, "average.png"),
	)
	if err != nil {
		return "", fmt.Errorf("error generating image average: %w", err)
	}

	avgImageFile, err := os.Open(avgImagePath)
	if err != nil {
		return "", fmt.Errorf("error opening %q: %w", avgImagePath, err)
	}

	posterImagePath, err := b.ip.OverlayImages(
		ctx,
		avgImageFile,
		bytes.NewReader(b.cfg.OverlayImage),
		filepath.Join(b.cfg.OutputDir, domain.PosterImageFileName),
	)
	if err != nil {
		return "", fmt.Errorf("image overlay error: %w", err)
	}
	if !b.cfg.SkipUpload {
		posterImageFileKey := filepath.Join(uploadDir, domain.PosterImageFileName)
		log.Infof("uploading %q as %q", posterImagePath, posterImageFileKey)
		posterImageFileKey, err = b.cs.UploadFile(ctx, posterImagePath, posterImageFileKey, "image/png", false)
		if err != nil {
			return "", fmt.Errorf("error uploading %q: %w", posterImageFileKey, err)
		}
	}

	return domain.PosterImageFileName, nil
}
