package videoprocessor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/mediatool"
)

func (p *Processor) CutVideo(ctx context.Context, inFile, outFile string, edit domain.VideoEdit) (string, error) {

	ext := filepath.Ext(outFile)
	progressFile := strings.Replace(outFile, ext, fmt.Sprintf("-progress%s", ext), 1)

	err := p.mt.CutVideo(ctx, inFile, progressFile, mediatool.Duration(edit.Start()), mediatool.Duration(edit.Duration()))
	if err != nil {
		return "", fmt.Errorf("error cutting video segment for %q: %v: %w", inFile, edit, err)
	}

	err = os.Rename(progressFile, outFile)
	if err != nil {
		return "", fmt.Errorf("error renaming %q to %q: %w", progressFile, outFile, err)
	}

	return outFile, nil
}
