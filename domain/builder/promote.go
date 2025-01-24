package builder

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (b *Builder) Promote(ctx context.Context, uploadDir string) error {

	if b.cfg.SkipUpload {
		log.Info("SkipUpload is true, skipping Promote step")
		return nil
	}

	r := strings.NewReader(uploadDir)
	_, err := b.cs.UploadReader(ctx, r, domain.CurrentBuildIDFileName, "text/plain", false)
	if err != nil {
		return fmt.Errorf("error uploading %q: %w", domain.CurrentBuildIDFileName, err)
	}
	return nil
}
