package builder

import (
	"context"
	"strings"
)

func (b *Builder) videoUploadStream(ctx context.Context, filePaths <-chan string, uploadDir string) <-chan string {

	// create upload fileKeys and send on 2-string channel
	uploadPaths := make(chan [2]string)
	go func() {
		defer close(uploadPaths)
		for filePath := range filePaths {
			fileKey := strings.Replace(filePath, b.cfg.OutputDir, uploadDir, 1)
			uploadPaths <- [2]string{filePath, fileKey}
		}
	}()

	return b.os.UploadFileStream(ctx, b.errs, uploadPaths, "video/webm", false)
}

// func (b *Builder) videoUploadStream(ctx context.Context, uploadDir string, clipFiles <-chan string) <-chan string {
// 	stream := make(chan string)

// 	var wg sync.WaitGroup

// 	cleanup := func() {
// 		wg.Wait()
// 		close(stream)
// 	}

// 	uploadVideo := func(filePath string) {
// 		defer wg.Done()

// 		var err error

// 		fileKey := strings.Replace(filePath, b.cfg.OutputDir, uploadDir, 1)

// 		log.Infof("uploading video %q as %q", filePath, fileKey)

// 		fileKey, err = b.up.UploadFile(ctx, filePath, fileKey, "video/webm", false)
// 		if err != nil {
// 			b.error(domain.ErrTypeUploadVideo, fmt.Errorf("error uploading video %q as %q; %w", filePath, fileKey, err))
// 			return
// 		}

// 		log.Infof("finished uploading video %q", fileKey)

// 		stream <- fileKey
// 	}

// 	go func() {
// 		defer cleanup()

// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			default:
// 				filePath, open := <-clipFiles

// 				if !open {
// 					return
// 				}

// 				wg.Add(1)
// 				go uploadVideo(filePath)
// 			}
// 		}
// 	}()

// 	return stream
// }
