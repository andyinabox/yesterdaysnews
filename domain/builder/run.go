package builder

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/streams"
)

func (b *Builder) Run(ctx context.Context) error {
	var err error

	defer func() {
		slog.Info("cleaning up...")
		removed, err := b.Cleanup(ctx, b.cfg.TotalBuildsToKeep)
		if err != nil {
			b.errs <- errorhandler.Err(domain.ErrTypeCleanup, fmt.Errorf("error during Cleanup phase: %w", err))
		}
		slog.Info("removed objects from object store", "count", len(removed))
	}()

	slog.Info("running setup...")
	err = b.Setup(ctx)
	if err != nil {
		return fmt.Errorf("error during Setup phase: %w", err)
	}

	// create additional build vars
	manifest := b.createManifest()

	slog.Info("building video clips...")
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

	if len(manifest.Files.Clips) == 0 {
		return errors.New("no clips were processed")
	}

	slog.Info("processed clips", "count", len(manifest.Files.Clips))

	// right now this step takes up too much memory to run in the serverless job
	// so leaving it disabled for now. might be better to use cli tool or something?
	slog.Info("skipping poster image generation...")
	// log.Info("generating poster image...")
	// // eventually it would be nice to not break the stream here
	// finishedClipsStream := streams.StringTransformStream(ctx, streams.StringStream(ctx, manifest.Files.Clips...), func(s string) string {
	// 	return filepath.Join(b.cfg.OutputDir, s)
	// })
	// extractedImagesStream := b.ExtractImages(ctx, finishedClipsStream)
	// posterImage, err := b.PosterImage(ctx, manifest.ID, streams.StringSlice(ctx, extractedImagesStream))
	// manifest.Files.PosterImageFile = posterImage
	// log.Infof("successfully uploaded %q", posterImage)

	slog.Info("building model...")
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
	slog.Info("successfully uploaded model", "key", modelFile)

	slog.Info("combining video files and generating subtitles...")
	clipFiles, err := filepath.Glob(filepath.Join(b.cfg.OutputDir, domain.ClipsDirName, "*.webm"))
	if err != nil {
		return fmt.Errorf("error getting video clip files: %w", err)
	}
	videoFile, subsFile, err := b.GenerateCombinedVideo(
		ctx,
		domain.ArchivePrefix,
		filepath.Join(b.cfg.OutputDir, modelFile),
		clipFiles,
	)
	if err != nil {
		return fmt.Errorf("error generating combined video: %w", err)
	}
	manifest.Files.VideoFile = videoFile
	manifest.Files.SubtitlesFile = subsFile

	slog.Info("uploading manifest...")
	manifestKey, err := b.Manifest(ctx, manifest.ID, manifest)
	if err != nil {
		return fmt.Errorf("error during Manifest phase: %w", err)
	}
	slog.Info("successfully uploaded manifest", "key", manifestKey)

	slog.Info("promoting uploaded files to current...")
	err = b.Promote(ctx, manifest.ID)
	if err != nil {
		return fmt.Errorf("error during Promote phase: %w", err)
	}
	slog.Info("promoted build to current", "buildID", manifest.ID)

	return nil
}
