package videoprocessor

import (
	"context"
	"math/rand/v2"
	"sync"

	"code.andydayton.com/andy/yesterdaysnews/domain"
	"code.andydayton.com/andy/yesterdaysnews/domain/errorhandler"
)

func (p *Processor) ShuffleClipsAndCombine(ctx context.Context, errs chan<- domain.Error, files []string, outFile string) (file string, err error) {

	files = p.validateVideoFiles(ctx, errs, files)

	// todo: make this seeded?
	rand.Shuffle(len(files), func(i, j int) {
		files[i], files[j] = files[j], files[i]
	})

	return p.mt.CombineVideos(ctx, files, outFile)
}

func (p *Processor) validateVideoFiles(ctx context.Context, errs chan<- domain.Error, files []string) (validFiles []string) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	validFiles = []string{}

	wg.Add(len(files))
	for _, fn := range files {
		go func() {
			defer wg.Done()

			err := p.mt.Validate(ctx, fn)
			if err != nil {
				errs <- errorhandler.Err(domain.ErrTypeValidateVideo, err)
				return
			}

			mu.Lock()
			validFiles = append(validFiles, fn)
			mu.Unlock()
		}()
	}
	wg.Wait()

	return
}
