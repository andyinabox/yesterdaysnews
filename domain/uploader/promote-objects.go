package uploader

import (
	"context"
	"errors"
)

func (u *Uploader) PromoteObjects(ctx context.Context, prefix string) (string, error) {

	// baseName := u.cfg.BucketNameBase
	// newBucketName := ""

	// baseExists, err := u.osclient.BucketExists(ctx, baseName)
	// if err != nil {
	// 	return "", err
	// }

	// // we need to rename the base bucket back to a suffixed one
	// if baseExists {
	// 	// fetching the manifest.json to get deployment ID
	// 	data, err := u.osclient.GetObject(ctx, baseName, "manifest.json")
	// 	if err != nil {
	// 		newBucketName = fmt.Sprintf("%s-%s", baseName, strconv.FormatInt(time.Now().Unix(), 10))
	// 		log.Errorf("no manifest found in base bucket, using name %q: %s", newBucketName, err)
	// 	} else {
	// 		manifest := Manifest{}
	// 		err = json.Unmarshal(data, &manifest)
	// 		if err != nil {
	// 			return "", fmt.Errorf("error unmarshaling manifest: %w", err)
	// 		}

	// 		newBucketName = fmt.Sprintf("%s-%s", baseName, manifest.ID)
	// 	}

	// 	log.Debugf("moving bucket %q to %q", baseName, newBucketName)
	// 	err = u.osclient.RenameBucket(ctx, baseName, newBucketName)
	// 	if err != nil {
	// 		return "", fmt.Errorf("error demoting base bucket to %q: %w", newBucketName, err)
	// 	}
	// }

	// log.Debugf("moving bucket %q to %q", containerName, baseName)
	// err = u.osclient.RenameBucket(ctx, containerName, baseName)
	// if err != nil {
	// 	return "", fmt.Errorf("error promoting bucket %s: %w", containerName, err)
	// }

	// return newBucketName, nil

	return "", errors.New("not implemented")
}
