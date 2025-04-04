package builder

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (b *Builder) Promote(ctx context.Context, buildId string) error {

	if b.cfg.SkipUpload {
		slog.Info("SkipUpload is true, skipping Promote step")
		return nil
	}

	r := strings.NewReader(buildId)
	_, err := b.cs.UploadReader(ctx, r, domain.CurrentBuildIDFileName, "text/plain", false)
	if err != nil {
		return fmt.Errorf("error uploading %q: %w", domain.CurrentBuildIDFileName, err)
	}
	return nil
}
