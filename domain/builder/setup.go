package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (b *Builder) Setup(ctx context.Context) error {
	dir := filepath.Join(b.cfg.OutputDir, domain.ClipsDirName)
	log.Infof("creating dir %q", dir)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating dir %q: %w", dir, err)
	}
	return nil
}
