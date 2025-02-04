package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"reflect"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

var inputGlob, outputFile string
var verbose bool

func init() {
	flag.StringVar(&inputGlob, "i", "dist/clips/*.png", "input file glob")
	flag.StringVar(&outputFile, "o", "dist/average.png", "output file")
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

const imageWidth = 1280
const imageHeight = 720

func main() {

	imagePaths, err := filepath.Glob(inputGlob)
	if err != nil {
		log.Fatal(err)
	}

	pixelTotals := make([][3]uint32, imageWidth*imageHeight)
	finalImage := image.NewNRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	for _, path := range imagePaths {
		log.Infof("process image %q", path)

		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			log.Fatal(err)
		}

		bounds := img.Bounds()

		if !reflect.DeepEqual(bounds, finalImage.Bounds()) {
			log.Errorf("image %q has incorrect bounds: %v", path, bounds)
			continue
		}

		i := 0
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				img1 := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
				pixelTotals[i] = [3]uint32{
					pixelTotals[i][0] + uint32(img1.R),
					pixelTotals[i][1] + uint32(img1.G),
					pixelTotals[i][2] + uint32(img1.B),
				}
				i++
			}
		}
	}

	log.Info("calculating averages")
	bounds := finalImage.Bounds()
	i := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.NRGBA{
				uint8(pixelTotals[i][0] / uint32(len(imagePaths))),
				uint8(pixelTotals[i][1] / uint32(len(imagePaths))),
				uint8(pixelTotals[i][2] / uint32(len(imagePaths))),
				255}
			finalImage.SetNRGBA(x, y, c)
			i++
		}
	}

	log.Info("saving final image")
	output, err := os.Create(outputFile)
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
