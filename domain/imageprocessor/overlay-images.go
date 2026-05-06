package imageprocessor

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log/slog"
	"os"
	"reflect"

	"code.andydayton.com/andy/yesterdaysnews/pkg/util"
)

func (p *Processor) OverlayImages(ctx context.Context, baseReader, overlayReader io.Reader, outputFilePath string) (string, error) {

	base, _, err := image.Decode(baseReader)
	if err != nil {
		return "", fmt.Errorf("error decoding base image: %w", err)
	}

	overlay, _, err := image.Decode(overlayReader)
	if err != nil {
		return "", fmt.Errorf("error decoding overlay image: %w", err)
	}

	// using base bounds here
	bounds := base.Bounds()

	if !reflect.DeepEqual(bounds, overlay.Bounds()) {
		return "", errors.New("base and overlay image are different sizes")
	}

	img := image.NewNRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			baseColor := color.NRGBAModel.Convert(base.At(x, y)).(color.NRGBA)
			overlayColor := color.NRGBAModel.Convert(overlay.At(x, y)).(color.NRGBA)

			c := color.NRGBA{
				util.LerpUint8(baseColor.R, overlayColor.R, overlayColor.A),
				util.LerpUint8(baseColor.G, overlayColor.G, overlayColor.A),
				util.LerpUint8(baseColor.B, overlayColor.B, overlayColor.A),
				255,
			}

			img.SetNRGBA(x, y, c)
		}
	}

	slog.Info("saving final image")
	output, err := os.Create(outputFilePath)
	if err != nil {
		return "", fmt.Errorf("error opening output file: %w", err)
	}
	defer output.Close()

	err = png.Encode(output, img)
	if err != nil {
		return "", fmt.Errorf("error encoding output file: %w", err)

	}

	return outputFilePath, nil
}
