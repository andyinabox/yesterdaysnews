package builder

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func (b *Builder) Setup(ctx context.Context) error {
	var err error

	err = b.yt.Setup(ctx)
	if err != nil {
		return err
	}

	if util.DoesFileExist(b.cfg.OutputDir) {
		slog.Info("removing contents of dir", "dir", b.cfg.OutputDir)
		err = util.RemoveContents(b.cfg.OutputDir)
		if err != nil {
			return fmt.Errorf("error removing dir %q: %w", b.cfg.OutputDir, err)
		}
	}

	clipsDir := filepath.Join(b.cfg.OutputDir, domain.ClipsDirName)

	slog.Info("creating dir", "dir", clipsDir)
	err = os.MkdirAll(clipsDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating dir %q: %w", clipsDir, err)
	}
	return nil
}
