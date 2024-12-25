package uploader

import (
	"context"
	"errors"
)

func (u *Uploader) Upload(ctx context.Context, dir string, manifest *Manifest) (containerName string, err error) {

	// manifest, err := u.CreateManifest(ctx, dir)
	// if err != nil {
	// 	return "", fmt.Errorf("error creating manifest: %w", err)
	// }

	// u.osclient.UploadFile()

	return "", errors.New("not implemented")
}
