package imageprocessor

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"reflect"
	"sync"

	"github.com/charmbracelet/log"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
)

func (p *Processor) AverageImages(ctx context.Context, errs chan<- domain.Error, paths []string, width, height int, outputFile string) (string, error) {

	var wg sync.WaitGroup
	var mu sync.Mutex

	finalImage := image.NewNRGBA(image.Rect(0, 0, width, height))
	finalImageBounds := finalImage.Bounds()
	pixelTotals := make([][3]uint32, width*height)
	completedCount := 0

	wg.Add(len(paths))
	for _, path := range paths {
		go func() {
			defer wg.Done()

			f, err := os.Open(path)
			if err != nil {
				errs <- errorhandler.Err(domain.ErrTypeTODO, err)
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				log.Fatal(err)
			}

			bounds := img.Bounds()

			if !reflect.DeepEqual(bounds, finalImageBounds) {
				errs <- errorhandler.Err(domain.ErrTypeTODO, fmt.Errorf("image %q has incorrect bounds: %v", path, bounds))
				return
			}

			i := 0
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					img1 := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)

					// not sure how this works ... maybe use a RWLock?
					// or maybe lock the entire loop?
					mu.Lock()
					pixelTotals[i] = [3]uint32{
						pixelTotals[i][0] + uint32(img1.R),
						pixelTotals[i][1] + uint32(img1.G),
						pixelTotals[i][2] + uint32(img1.B),
					}
					mu.Unlock()

					i++
				}
			}

			// this could also be a channel with file paths but not sure if that's necessary
			mu.Lock()
			completedCount++
			mu.Unlock()

		}()
	}

	wg.Wait()

	i := 0
	for y := finalImageBounds.Min.Y; y < finalImageBounds.Max.Y; y++ {
		for x := finalImageBounds.Min.X; x < finalImageBounds.Max.X; x++ {
			c := color.NRGBA{
				uint8(pixelTotals[i][0] / uint32(completedCount)),
				uint8(pixelTotals[i][1] / uint32(completedCount)),
				uint8(pixelTotals[i][2] / uint32(completedCount)),
				255}
			finalImage.SetNRGBA(x, y, c)
			i++
		}
	}

	log.Debug("saving final image")
	output, err := os.Create(outputFile)
	if err != nil {
		return "", fmt.Errorf("error opening output image file: %w", err)
	}

	if err := png.Encode(output, finalImage); err != nil {
		output.Close()
		return "", fmt.Errorf("error encoding output image file: %w", err)
	}

	if err := output.Close(); err != nil {
		return "", fmt.Errorf("error closing output image file: %w", err)
	}

	return outputFile, nil
}
