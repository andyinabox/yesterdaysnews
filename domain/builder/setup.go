package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func (b *Builder) Setup(ctx context.Context) error {
	var err error

	if util.DoesFileExist(b.cfg.OutputDir) {
		log.Infof("removing contents of dir %q", b.cfg.OutputDir)
		err = util.RemoveContents(b.cfg.OutputDir)
		if err != nil {
			return fmt.Errorf("error removing dir %q: %w", b.cfg.OutputDir, err)
		}
	}

	clipsDir := filepath.Join(b.cfg.OutputDir, domain.ClipsDirName)

	log.Infof("creating dir %q", clipsDir)
	err = os.MkdirAll(clipsDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating dir %q: %w", clipsDir, err)
	}
	return nil
}
