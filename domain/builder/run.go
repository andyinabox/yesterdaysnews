package builder

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
)

func (b *Builder) Run(ctx context.Context) error {
	var err error

	defer func() {
		log.Info("cleaning up...")
		removed, err := b.Cleanup(ctx, b.cfg.TotalBuildsToKeep)
		if err != nil {
			b.errs <- errorhandler.Err(domain.ErrTypeCleanup, fmt.Errorf("error during Cleanup phase: %w", err))
		}
		log.Infof("removed %d objects from object store", len(removed))
	}()

	log.Info("running setup...")
	err = b.Setup(ctx)
	if err != nil {
		return fmt.Errorf("error during Setup phase: %w", err)
	}

	// create additional build vars
	manifest := b.createManifest()

	log.Info("building video clips...")
	manifest.Files.Clips, err = b.VideoClips(ctx, manifest.ContentDate, manifest.ID)
	if err != nil {
		return fmt.Errorf("error during VideoClips phase: %w", err)
	}
	log.Infof("processed %d clips", len(manifest.Files.Clips))

	log.Info("building model...")
	modelFile, err := b.Model(ctx, manifest.ID, []domain.Corpus{
		{
			Type:     domain.CorpusTypeVTT,
			FileGlob: filepath.Join(b.cfg.OutputDir, "*.vtt"),
			Weight:   b.cfg.CaptionNewsCorpusWeight,
		},
		{
			Type:    domain.CorpusTypeString,
			Content: b.cfg.HospitalCorpus,
			Weight:  b.cfg.CaptionHospitalCorpusWeight,
		},
	})
	if err != nil {
		return fmt.Errorf("error during Model phase: %w", err)
	}
	manifest.Files.ModelFile = modelFile
	log.Infof("successfully uploaded %q", modelFile)

	log.Info("uploading manifest...")
	manifestKey, err := b.Manifest(ctx, manifest.ID, manifest)
	if err != nil {
		return fmt.Errorf("error during Manifest phase: %w", err)
	}
	log.Infof("successfully uploaded %q", manifestKey)

	log.Info("promoting uploaded files to current...")
	err = b.Promote(ctx, manifest.ID)
	if err != nil {
		return fmt.Errorf("error during Promote phase: %w", err)
	}
	log.Infof("promoted %q to current", manifest.ID)

	return nil
}
