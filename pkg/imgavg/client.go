package imgavg

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

const batchSize = 255

type Client struct {
	width      int
	height     int
	sum        pixelSum
	dir        string
	errHandler func(error)
}

func New(dir string, width, height int, errHandler func(error)) *Client {
	return &Client{
		width:      width,
		height:     height,
		sum:        make(pixelSum, width*height),
		dir:        dir,
		errHandler: errHandler,
	}
}

func (c *Client) Run(paths []string, outFile string) (string, error) {
	return c.run(paths, outFile, 0)
}

func (c *Client) run(paths []string, outFile string, count int) (string, error) {
	var err error

	// single batch
	if len(paths) <= batchSize {
		result, err := c.avg(paths, count)
		if err != nil {
			return "", fmt.Errorf("error running single batch %d: %w", count, err)
		}
		err = os.Rename(result, outFile)
		if err != nil {
			return "", fmt.Errorf("error moving completed image from %q to %q: %w", paths[0], outFile, err)
		}
		return outFile, nil
	}

	// break into batches that will generate individual files,
	// then pass those files into this method recursivley
	batchCount := (len(paths) / batchSize) + 1
	batchPaths := make([]string, batchCount)

	log.Debugf("breaking into %d batches", batchCount)
	var start, end int
	for i := 0; i < batchCount; i++ {
		start = batchSize * i
		end = start + batchSize
		if len(paths) < end {
			end = len(paths)
		}

		count++
		batchPaths[i], err = c.avg(batchPaths[start:end], count)
		if err != nil {
			return "", fmt.Errorf("error in level %d: %w", count, err)
		}
	}

	count++
	log.Debugf("finished creating intermediary images, now running batch %d", count)
	return c.run(batchPaths, outFile, count)
}

func (c *Client) avg(paths []string, count int) (string, error) {

	c.sum.Clear()

	completed := 0

	for _, path := range paths {

		log.Debugf("collect pixels for %q", path)

		// open file
		f, err := os.Open(path)
		if err != nil {
			c.errHandler(fmt.Errorf("error opening image %q: %w", path, err))
			continue
		}
		defer f.Close()

		// decode image
		img, _, err := image.Decode(f)
		if err != nil {
			c.errHandler(fmt.Errorf("error decoding image %q: %w", path, err))
			continue
		}

		// get image bounds
		b := img.Bounds()

		if b.Max.X*b.Max.Y != len(c.sum) {
			c.errHandler(fmt.Errorf("image %q has incorrect dimensions: %d x %d", path, b.Max.X, b.Max.Y))
			continue
		}

		// iterate through image coords and add to pixels array
		var pix int
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				img := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
				c.sum[pix][0] = c.sum[pix][0] + uint16(img.R)
				c.sum[pix][1] = c.sum[pix][1] + uint16(img.G)
				c.sum[pix][2] = c.sum[pix][2] + uint16(img.B)
				pix++
			}
		}

		completed++
	}

	img := image.NewNRGBA(image.Rect(0, 0, c.width, c.height))
	bounds := img.Bounds()

	var pix int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.NRGBA{
				uint8(c.sum[pix][0] / uint16(completed)),
				uint8(c.sum[pix][1] / uint16(completed)),
				uint8(c.sum[pix][2] / uint16(completed)),
				255}
			img.SetNRGBA(x, y, c)
			pix++
		}
	}

	outFile := filepath.Join(c.dir, fmt.Sprintf("%d.png", count))
	output, err := os.Create(outFile)
	if err != nil {
		return "", fmt.Errorf("error opening output image file: %w", err)
	}
	defer output.Close()

	err = png.Encode(output, img)
	if err != nil {
		return "", fmt.Errorf("error encoding output image file: %w", err)
	}

	return outFile, nil
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
