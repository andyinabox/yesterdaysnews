package builder

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func (b *Builder) Run(ctx context.Context) error {

	now := time.Now()
	manifest := &domain.Manifest{
		Date: now,
		ID:   util.Timestamp(now),
		Files: domain.ManifestFiles{
			Clips: []string{},
		},
	}

	yesterday := util.Yesterday()
	uploadDir := util.Timestamp(yesterday)

	log.Info("building video clips")
	manifest.Files.Clips = b.BuildVideoClips(ctx, yesterday, uploadDir)

	log.Info("building model")
	modelFile, err := b.BuildModel(ctx, uploadDir, []domain.Corpus{
		{
			Type:     domain.CorpusTypeVTT,
			FileGlob: "download/*.vtt",
			Weight:   1,
		},
		{
			Type:     domain.CorpusTypeText,
			FileGlob: "hospital.txt",
			Weight:   3,
		},
	})
	if err != nil {
		return err
	}
	manifest.Files.ModelFile = modelFile

	log.Info("uploading manifest")
	err = b.UploadManifest(ctx, uploadDir, manifest)
	if err != nil {
		return err
	}

	return nil
}
