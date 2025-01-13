package builder

import (
	"context"
	"fmt"
	"strings"
)

const currentFileName = "current.txt"

func (b *Builder) Promote(ctx context.Context, uploadDir string) error {
	r := strings.NewReader(uploadDir)
	_, err := b.cs.UploadReader(ctx, r, currentFileName, "text/plain", false)
	if err != nil {
		return fmt.Errorf("error uploading %q: %w", currentFileName, err)
	}
	return nil
}

// func (b *Builder) Promote(ctx context.Context, uploadDir string) (string, error) {

// 	// first move the files from the primary prefix to a different prefiex
// 	demotedName := b.getCurrentManifestName(ctx)
// 	log.Infof("demoting current objects to %q", demotedName)
// 	fileKeys, err := b.moveObjects(ctx, b.cfg.ObjectStorePrimaryDir, demotedName)
// 	if err != nil {
// 		return "", fmt.Errorf("error moving objects from %q to %q: %w", b.cfg.ObjectStorePrimaryDir, demotedName, err)
// 	}
// 	log.Infof("successfully moved %d objects from %q to %q", len(fileKeys), b.cfg.ObjectStorePrimaryDir, demotedName)

// 	log.Infof("promoting %q", uploadDir)
// 	fileKeys, err = b.moveObjects(ctx, uploadDir, b.cfg.ObjectStorePrimaryDir)
// 	if err != nil {
// 		return "", fmt.Errorf("error moving objects from %q to %q: %w", uploadDir, b.cfg.ObjectStorePrimaryDir, err)
// 	}
// 	log.Infof("successfully moved %d objects from %q to %q", len(fileKeys), uploadDir, b.cfg.ObjectStorePrimaryDir)

// 	return demotedName, nil
// }

// func (b *Builder) moveObjects(ctx context.Context, fromPrefix, toPrefix string) ([]string, error) {
// 	log.Debugf("move objects %q, %q", fromPrefix, toPrefix)

// 	// get slice of keys from list command
// 	fileKeys, err := b.cs.ListObjectsWithPrefix(ctx, fromPrefix)
// 	if err != nil {
// 		return nil, fmt.Errorf("error listing objects with %q prefix: %w", fromPrefix, err)
// 	}

// 	// convert to string stream with throttling
// 	fileKeyStream := streams.StringStreamThrottled(ctx, time.Millisecond, fileKeys...)

// 	// convert to [2]string stream for renaming
// 	movePathStream := streams.StringStreamTo2StringSliceStream(ctx, fileKeyStream, func(s string) [2]string {
// 		return [2]string{s, strings.Replace(s, fromPrefix, toPrefix, 1)}
// 	})

// 	// move objects
// 	movedObjectsStream := b.cs.MoveObjectStream(ctx, b.errs, movePathStream)

// 	// return a string slice
// 	return streams.StringSlice(ctx, movedObjectsStream), nil
// }

// // getCurrentManifestName attempts to get the build ID from the manifest in
// // the primary object store prefix/dir. If that fails default to a current timestamp
// func (b *Builder) getCurrentManifestName(ctx context.Context) (prefix string) {
// 	var err error
// 	prefix = util.Timestamp(time.Now())

// 	currentManifestFileKey := filepath.Join(b.cfg.ObjectStorePrimaryDir, "manifest.json")
// 	data, err := b.cs.GetObject(ctx, currentManifestFileKey)
// 	if err != nil {
// 		b.eh.Add(domain.ErrTypeGetCurrentManifestPrefix, fmt.Errorf("error downloading %q: %w", currentManifestFileKey, err))
// 		return
// 	}

// 	manifest := domain.Manifest{}
// 	err = json.Unmarshal(data, &manifest)
// 	if err != nil {
// 		b.eh.Add(domain.ErrTypeGetCurrentManifestPrefix, fmt.Errorf("error unmarshaling current manifest: %w", err))
// 		return
// 	}

// 	if manifest.ID != "" {
// 		prefix = manifest.ID
// 	}

// 	return
// }
