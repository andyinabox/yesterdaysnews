package builder

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (b *Builder) PosterImage(ctx context.Context, uploadDir string, paths <-chan string) (string, error) {

	avgImagePath, err := b.ip.AverageImagesStream(
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

	posteImageFileKey := filepath.Join(uploadDir, domain.PosterImageFileName)
	log.Infof("uploading %q as %q", posterImagePath, posteImageFileKey)
	posteImageFileKey, err = b.cs.UploadFile(ctx, posterImagePath, posteImageFileKey, "image/png", false)
	if err != nil {
		return "", fmt.Errorf("error uploading %q: %w", posteImageFileKey, err)
	}
	return domain.PosterImageFileName, nil
}
