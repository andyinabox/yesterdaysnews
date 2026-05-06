package builder

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"code.andydayton.com/andy/yesterdaysnews/domain"
	"code.andydayton.com/andy/yesterdaysnews/domain/captiongenerator"
	"code.andydayton.com/andy/yesterdaysnews/domain/captionschain"
	"code.andydayton.com/andy/yesterdaysnews/pkg/markov"
	"code.andydayton.com/andy/yesterdaysnews/pkg/util"
)

func (b *Builder) GenerateCombinedVideo(ctx context.Context, uploadDir, modelFile string, clips []string) (string, string, error) {

	videoFile := filepath.Join(b.cfg.OutputDir, domain.VideoFileName)

	if util.DoesFileExist(videoFile) {
		slog.Info("video file already exists, deleting", "file", videoFile)
		err := os.Remove(videoFile)
		if err != nil {
			return "", "", fmt.Errorf("error removing video file: %w", err)
		}
	}

	slog.Debug("combining video files", "count", len(clips), "file", videoFile)
	videoFile, err := b.vp.ShuffleClipsAndCombine(ctx, b.errs, clips, videoFile)
	if err != nil {
		return "", "", fmt.Errorf("error combining video clips: %w", err)
	}

	videoLength, err := b.vp.GetVideoLength(ctx, videoFile)
	if err != nil {
		return "", "", fmt.Errorf("error getting video length: %w", err)
	}

	chain, err := b.buildMarkovChain(modelFile)
	if err != nil {
		return "", "", fmt.Errorf("error building markov chain: %w", err)
	}

	cg := captiongenerator.New(chain, &captiongenerator.Config{
		MinCaptionLength: 7,
		MaxCaptionLength: 15,
	})

	slog.Debug("generating subtitles", "videoLength", videoLength)
	subs, err := cg.Subtitles(videoLength, b.cfg.CaptionMinDuration, b.cfg.CaptionMaxDuration)
	if err != nil {
		return "", "", fmt.Errorf("error generating subtitles: %w", err)
	}
	subsFile := filepath.Join(b.cfg.OutputDir, domain.SubsFileName)
	err = os.WriteFile(subsFile, subs, os.ModePerm)
	if err != nil {
		return "", "", fmt.Errorf("error writing subtitles file: %w", err)
	}

	if b.cfg.SkipUpload {
		return videoFile, subsFile, nil
	}

	dateStr := util.DateString(util.Yesterday())
	videoKey := filepath.Join(uploadDir, dateStr+".mp4")
	subsKey := filepath.Join(uploadDir, dateStr+".srt")

	slog.Info("uploading combined video file (this may take a while)", "file", videoFile, "key", videoKey)
	videoKey, err = b.cs.UploadFile(
		ctx,
		videoFile,
		videoKey,
		"video/mp4",
		true,
	)
	if err != nil {
		return "", "", fmt.Errorf("error uploading video file: %w", err)
	}

	slog.Info("uploading combined subs file", "file", subsFile, "key", subsKey)
	subsKey, err = b.cs.UploadFile(
		ctx,
		subsFile,
		subsKey,
		"text/plain",
		false,
	)
	if err != nil {
		return "", "", fmt.Errorf("error uploading subs file: %w", err)
	}

	return videoKey, subsKey, nil
}

func (b *Builder) buildMarkovChain(modelFile string) (markov.Chain, error) {

	modelData, err := os.ReadFile(modelFile)
	if err != nil {
		return nil, err
	}

	chain := captionschain.New(b.cfg.CaptionPrefixLength)
	err = chain.Load(modelData)
	if err != nil {
		return nil, err
	}

	return chain, nil
}
