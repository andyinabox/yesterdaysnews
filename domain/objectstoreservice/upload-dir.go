package objectstoreservice

// var ErrContainerDoesNotExist = errors.New("container does not exist")

// func (s *Service) UploadDir(ctx context.Context, dir string) (string, error) {

// 	var wg sync.WaitGroup

// 	manifestBytes, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
// 	if err != nil {
// 		return "", fmt.Errorf("error reading manifest: %w", err)
// 	}

// 	manifest := domain.Manifest{}
// 	err = json.Unmarshal(manifestBytes, &manifest)
// 	if err != nil {
// 		return "", fmt.Errorf("error decoding manifest: %w", err)
// 	}

// 	deployDir := manifest.ID

// 	exists, err := u.osclient.ContainerExists(ctx, u.cfg.ContainerName)
// 	if err != nil {
// 		return "", err
// 	}
// 	if !exists {
// 		return "", ErrContainerDoesNotExist
// 	}

// 	// upload file func
// 	uploadFile := func(path, contentType string, multipart bool) {
// 		defer wg.Done()

// 		key := filepath.Join(deployDir, path)
// 		filePath := filepath.Join(dir, path)

// 		log.Debugf("begin uploading file %q as %q", filePath, key)

// 		file, err := os.Open(filePath)
// 		if err != nil {
// 			log.Errorf("error opening file %q: %s", filePath, err)
// 		}

// 		_, err = u.osclient.UploadFile(
// 			ctx,
// 			u.cfg.ContainerName,
// 			key,
// 			file,
// 			contentType,
// 			multipart,
// 		)
// 		if err != nil {
// 			log.Errorf("error uploading file %q: %s", key, err)
// 		}

// 		log.Debugf("finished uploading %q", key)
// 	}

// 	// upload model file
// 	wg.Add(1)
// 	go uploadFile(manifest.Files.ModelFile, "application/json", false)

// 	// upload video file
// 	// wg.Add(1)
// 	// go uploadFile(manifest.Files.VideoFile, "video/mp4", true)

// 	// upload individual clips
// 	for _, clipPath := range manifest.Files.Clips {
// 		wg.Add(1)
// 		go uploadFile(clipPath, "video/webm", false)
// 	}

// 	wg.Wait()

// 	manifestPath := filepath.Join(deployDir, "manifest.json")
// 	log.Infof("uploading manifest to %s", manifestPath)
// 	_, err = u.osclient.UploadFile(ctx, u.cfg.ContainerName, manifestPath, bytes.NewReader(manifestBytes), "application/json", false)
// 	if err != nil {
// 		return "", fmt.Errorf("error uploading manifest file: %w", err)
// 	}

// 	return deployDir, nil
// }
