package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"
	"reflect"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

var baseImagePath, overlayImagePath, outputImagePath string
var verbose bool

func init() {
	flag.StringVar(&baseImagePath, "base", "dist/average.png", "base image on bottom layer")
	flag.StringVar(&overlayImagePath, "overlay", "app/builder/overlay.png", "overlay image on top layer")
	flag.StringVar(&outputImagePath, "o", "dist/poster.png", "output image location")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	err := godotenv.Load()
	if err != nil {
		log.Warnf("error getting .env: %s", err)
	}
}

func LerpUint8(a, b, t uint8) uint8 {
	return uint8(uint16(a) + (uint16(t) * (uint16(b) - uint16(a)) / 255))
}

func main() {

	baseImgFile, err := os.Open(baseImagePath)
	if err != nil {
		log.Debug(err)
	}
	baseImg, _, err := image.Decode(baseImgFile)
	if err != nil {
		log.Debug(err)
	}
	baseImgBounds := baseImg.Bounds()

	overlayImageFile, err := os.Open(overlayImagePath)
	if err != nil {
		log.Debug(err)
	}
	overlayImg, _, err := image.Decode(overlayImageFile)
	if err != nil {
		log.Debug(err)
	}

	if !reflect.DeepEqual(baseImgBounds, overlayImg.Bounds()) {
		log.Fatal("base and overlay image are different sizes")
	}

	finalImage := image.NewNRGBA(image.Rect(0, 0, baseImgBounds.Max.X, baseImgBounds.Max.Y))

	for y := baseImgBounds.Min.Y; y < baseImgBounds.Max.Y; y++ {
		for x := baseImgBounds.Min.X; x < baseImgBounds.Max.X; x++ {
			baseImgColor := color.NRGBAModel.Convert(baseImg.At(x, y)).(color.NRGBA)
			overlayImgColor := color.NRGBAModel.Convert(overlayImg.At(x, y)).(color.NRGBA)

			c := color.NRGBA{
				LerpUint8(baseImgColor.R, overlayImgColor.R, overlayImgColor.A),
				LerpUint8(baseImgColor.G, overlayImgColor.G, overlayImgColor.A),
				LerpUint8(baseImgColor.B, overlayImgColor.B, overlayImgColor.A),
				255,
			}

			finalImage.SetNRGBA(x, y, c)
		}
	}

	log.Info("saving final image")
	output, err := os.Create(outputImagePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(output, finalImage); err != nil {
		output.Close()
		log.Fatal(err)
	}

	if err := output.Close(); err != nil {
		log.Fatal(err)
	}
	log.Info("done")

}
