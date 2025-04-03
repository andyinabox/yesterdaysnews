package builder

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/captiongenerator"
	"gitlab.com/andyinabox/yesterdaysnews/domain/captionschain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/markov"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

const minCaptionDuration = 3.0
const maxCaptionDuration = 7.0

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
	videoFile, err := b.vp.ShuffleClipsAndCombine(ctx, clips, videoFile)
	if err != nil {
		return "", "", fmt.Errorf("error combining video clips: %w", err)
	}

	videoLength, err := b.vp.GetVideoLength(ctx, videoFile)
	if err != nil {
		return "", "", fmt.Errorf("error getting video length: %w", err)
	}

	chain, err := b.buildMarkovChain(ctx, modelFile)
	if err != nil {
		return "", "", fmt.Errorf("error building markov chain: %w", err)
	}

	cg := captiongenerator.New(chain, &captiongenerator.Config{
		MinCaptionLength: 7,
		MaxCaptionLength: 15,
	})

	slog.Debug("generating subtitles", "videoLength", videoLength)
	subs, err := cg.Subtitles(videoLength, minCaptionDuration, maxCaptionDuration)
	if err != nil {
		return "", "", fmt.Errorf("error generating subtitles: %w", err)
	}
	subsFile := filepath.Join(b.cfg.OutputDir, domain.SubsFileName)
	err = os.WriteFile(subsFile, subs, os.ModePerm)
	if err != nil {
		return "", "", fmt.Errorf("error writing subtitles file: %w", err)
	}

	return videoFile, "", nil
}

func (b *Builder) buildMarkovChain(ctx context.Context, modelFile string) (markov.Chain, error) {

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
