package videoprocessor

import (
	"context"
	"math/rand/v2"
)

func (p *Processor) ShuffleClipsAndCombine(ctx context.Context, files []string, outFile string) (file string, err error) {

	// todo: make this seeded?
	rand.Shuffle(len(files), func(i, j int) {
		files[i], files[j] = files[j], files[i]
	})

	return p.mt.CombineVideos(ctx, files, outFile)
}
