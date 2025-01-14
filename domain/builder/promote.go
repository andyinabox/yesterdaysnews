package builder

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (b *Builder) Promote(ctx context.Context, uploadDir string) error {
	r := strings.NewReader(uploadDir)
	_, err := b.cs.UploadReader(ctx, r, domain.CurrentBuildIDFileName, "text/plain", false)
	if err != nil {
		return fmt.Errorf("error uploading %q: %w", domain.CurrentBuildIDFileName, err)
	}
	return nil
}
