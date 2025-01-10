package builder

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func (b *Builder) Run(ctx context.Context) error {
	var err error

	log.Info("running setup...")
	err = b.Setup(ctx)
	if err != nil {
		return fmt.Errorf("error during Setup phase: %w", err)
	}

	// create additional build vars
	manifest := b.createManifest()
	yesterday := util.Yesterday()
	uploadDir := util.Timestamp(yesterday)

	log.Info("building video clips...")
	manifest.Files.Clips, err = b.VideoClips(ctx, yesterday, uploadDir)
	if err != nil {
		return fmt.Errorf("error during VideoClips phase: %w", err)
	}
	log.Infof("processed %d clips", len(manifest.Files.Clips))

	log.Info("building model...")
	modelFile, err := b.Model(ctx, uploadDir, []domain.Corpus{
		{
			Type:     domain.CorpusTypeVTT,
			FileGlob: filepath.Join(b.cfg.OutputDir, "*.vtt"),
			Weight:   1,
		},
		{
			Type:     domain.CorpusTypeText,
			FileGlob: "hospital.txt",
			Weight:   3,
		},
	})
	if err != nil {
		return fmt.Errorf("error during Model phase: %w", err)
	}
	manifest.Files.ModelFile = modelFile
	log.Infof("successfully uploaded %q", modelFile)

	log.Info("uploading manifest...")
	manifestKey, err := b.Manifest(ctx, uploadDir, manifest)
	if err != nil {
		return fmt.Errorf("error during Manifest phase: %w", err)
	}
	log.Infof("successfully uploaded %q", manifestKey)

	log.Info("promoting uploaded files to current...")
	demoted, err := b.Promote(ctx, uploadDir)
	if err != nil {
		return fmt.Errorf("error during Promote phase: %w", err)
	}
	log.Infof("promoted new files and demoted %q", demoted)

	log.Info("cleaning up...")
	removed, err := b.Cleanup(ctx, demoted)
	if err != nil {
		return fmt.Errorf("error during Cleanup phase: %w", err)
	}
	log.Infof("removed %d objects from object store", len(removed))

	return nil
}

func (b *Builder) createManifest() *domain.Manifest {
	now := time.Now()
	return &domain.Manifest{
		Date: now,
		ID:   util.Timestamp(now),
		Files: domain.ManifestFiles{
			Clips: []string{},
		},
	}
}
