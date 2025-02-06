package builder

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/streams"
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
	downloadPathStream := b.DownloadVideos(ctx, manifest.ContentDate)
	clipPathsStream := b.CutVideos(ctx, downloadPathStream)

	// skip uploading files
	if b.cfg.SkipUpload {
		clipPathsStream = streams.StringTransformStream(ctx, clipPathsStream, func(s string) string {
			return strings.TrimPrefix(s, b.cfg.OutputDir+"/")
		})
		manifest.Files.Clips = streams.StringSlice(ctx, clipPathsStream)

		// upload files
	} else {
		uploadPathsStream := b.UploadVideos(ctx, manifest.ID, clipPathsStream)
		manifest.Files.Clips = streams.StringSlice(ctx, uploadPathsStream)
	}
	log.Infof("processed %d clips", len(manifest.Files.Clips))

	// right now this step takes up too much memory to run in the serverless job
	// so leaving it disabled for now. might be better to use cli tool or something?
	log.Info("skipping poster image generation...")
	// log.Info("generating poster image...")
	// // eventually it would be nice to not break the stream here
	// finishedClipsStream := streams.StringTransformStream(ctx, streams.StringStream(ctx, manifest.Files.Clips...), func(s string) string {
	// 	return filepath.Join(b.cfg.OutputDir, s)
	// })
	// extractedImagesStream := b.ExtractImages(ctx, finishedClipsStream)
	// posterImage, err := b.PosterImage(ctx, manifest.ID, streams.StringSlice(ctx, extractedImagesStream))
	// manifest.Files.PosterImageFile = posterImage
	// log.Infof("successfully uploaded %q", posterImage)

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
