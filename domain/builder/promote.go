package builder

import (
	"context"
	"fmt"
	"strings"
)

const currentFileName = "current.txt"

func (b *Builder) Promote(ctx context.Context, uploadDir string) error {
	r := strings.NewReader(uploadDir)
	_, err := b.cs.UploadReader(ctx, r, currentFileName, "text/plain", false)
	if err != nil {
		return fmt.Errorf("error uploading %q: %w", currentFileName, err)
	}
	return nil
}
