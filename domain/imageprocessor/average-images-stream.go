package imageprocessor

// func (p *Processor) AverageImagesStream(ctx context.Context, errs chan<- domain.Error, paths <-chan string, width, height int, outputFile string) (string, error) {

// 	bounds := image.Rect(0, 0, width, height)
// 	img := image.NewNRGBA(bounds)

// 	completedPaths, pixels := p.collectPixels(ctx, errs, paths, bounds)

// 	count := len(completedPaths)

// 	var i int
// 	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
// 		for x := bounds.Min.X; x < bounds.Max.X; x++ {
// 			c := color.NRGBA{
// 				uint8(pixels[i][0] / uint32(count)),
// 				uint8(pixels[i][1] / uint32(count)),
// 				uint8(pixels[i][2] / uint32(count)),
// 				255}
// 			img.SetNRGBA(x, y, c)
// 			i++
// 		}
// 	}

// 	log.Debug("saving final image")
// 	output, err := os.Create(outputFile)
// 	if err != nil {
// 		return "", fmt.Errorf("error opening output image file: %w", err)
// 	}
// 	defer output.Close()

// 	err = png.Encode(output, img)
// 	if err != nil {
// 		return "", fmt.Errorf("error encoding output image file: %w", err)
// 	}

// 	return outputFile, nil
// }

// func (p *Processor) collectPixels(ctx context.Context, errs chan<- domain.Error, paths <-chan string, bounds image.Rectangle) ([]string, [][3]uint32) {

// 	var wg sync.WaitGroup
// 	var mu sync.Mutex

// 	pixels := make([][3]uint32, bounds.Max.X*bounds.Max.Y)
// 	stream := make(chan string)

// 	cleanup := func() {
// 		wg.Wait()
// 		close(stream)
// 	}

// 	collect := func(path string) {
// 		defer wg.Done()

// 		log.Debugf("collect pixels for %q", path)

// 		// open file
// 		f, err := os.Open(path)
// 		if err != nil {
// 			errs <- errorhandler.Err(domain.ErrTypeAverageImage, fmt.Errorf("error opening image %q: %w", path, err))
// 			return
// 		}
// 		defer f.Close()

// 		// decode image
// 		img, _, err := image.Decode(f)
// 		if err != nil {
// 			errs <- errorhandler.Err(domain.ErrTypeAverageImage, fmt.Errorf("error decoding image %q: %w", path, err))
// 			return
// 		}

// 		// get image bounds
// 		b := img.Bounds()

// 		if !reflect.DeepEqual(b, bounds) {
// 			errs <- errorhandler.Err(domain.ErrTypeAverageImage, fmt.Errorf("image %q has incorrect bounds: %v", path, b))
// 			return
// 		}

// 		// it works MUCH better when locking for the entire nested loop
// 		mu.Lock()
// 		// iterate through image coords and add to pixels array
// 		var i int
// 		for y := b.Min.Y; y < b.Max.Y; y++ {
// 			for x := b.Min.X; x < b.Max.X; x++ {
// 				img1 := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
// 				pixels[i] = [3]uint32{
// 					pixels[i][0] + uint32(img1.R),
// 					pixels[i][1] + uint32(img1.G),
// 					pixels[i][2] + uint32(img1.B),
// 				}
// 				i++
// 			}
// 		}
// 		mu.Unlock()

// 		log.Debugf("successfully collected pixels for %q", path)
// 		stream <- path
// 	}

// 	// our little concurrency nugget <3
// 	go func() {
// 		defer cleanup()
// 		for path := range paths {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			default:
// 				wg.Add(1)
// 				go collect(path)
// 			}
// 		}
// 	}()

// 	// block until stream is closed and then return completed paths and pixels
// 	completedPaths := streams.StringSlice(ctx, stream)

// 	return completedPaths, pixels
// }
