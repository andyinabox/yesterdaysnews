package imageprocessor

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"reflect"

	"github.com/charmbracelet/log"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/errorhandler"
)

func (p *Processor) AverageImages(ctx context.Context, errs chan<- domain.Error, paths []string, width, height int, outputFile string) (string, error) {
	var count int

	bounds := image.Rect(0, 0, width, height)
	pixels := make([][3]uint32, bounds.Max.X*bounds.Max.Y)

	for _, path := range paths {

		log.Debugf("collect pixels for %q", path)

		// open file
		f, err := os.Open(path)
		if err != nil {
			errs <- errorhandler.Err(domain.ErrTypeAverageImage, fmt.Errorf("error opening image %q: %w", path, err))
			continue
		}
		defer f.Close()

		// decode image
		img, _, err := image.Decode(f)
		if err != nil {
			errs <- errorhandler.Err(domain.ErrTypeAverageImage, fmt.Errorf("error decoding image %q: %w", path, err))
			continue
		}

		// get image bounds
		b := img.Bounds()

		if !reflect.DeepEqual(b, bounds) {
			errs <- errorhandler.Err(domain.ErrTypeAverageImage, fmt.Errorf("image %q has incorrect bounds: %v", path, b))
			continue
		}

		// iterate through image coords and add to pixels array
		var i int
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				img1 := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
				pixels[i] = [3]uint32{
					pixels[i][0] + uint32(img1.R),
					pixels[i][1] + uint32(img1.G),
					pixels[i][2] + uint32(img1.B),
				}
				i++
			}
		}

		log.Debugf("successfully collected pixels for %q", path)

		count++
	}

	img := image.NewNRGBA(bounds)

	var i int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.NRGBA{
				uint8(pixels[i][0] / uint32(count)),
				uint8(pixels[i][1] / uint32(count)),
				uint8(pixels[i][2] / uint32(count)),
				255}
			img.SetNRGBA(x, y, c)
			i++
		}
	}

	log.Debug("saving final image")
	output, err := os.Create(outputFile)
	if err != nil {
		return "", fmt.Errorf("error opening output image file: %w", err)
	}
	defer output.Close()

	err = png.Encode(output, img)
	if err != nil {
		return "", fmt.Errorf("error encoding output image file: %w", err)
	}

	return outputFile, nil
}
