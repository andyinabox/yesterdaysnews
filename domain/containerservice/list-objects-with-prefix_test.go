//go:build objectstoretest
// +build objectstoretest

package containerservice

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func TestListObjectsWithPrefix(t *testing.T) {
	s := New(&Config{
		S3Endpoint:    os.Getenv("YN_S3_ENDPOINT"),
		S3AccessKey:   os.Getenv("YN_S3_ACCESS_KEY"),
		S3SecretKey:   os.Getenv("YN_S3_SECRET_ACCESS_KEY"),
		S3Region:      os.Getenv("YN_S3_REGION"),
		ContainerName: "yesterdaysnewstest",
	})

	ctx := context.Background()
	timestamp := util.Timestamp(time.Now())

	filePath := "../../test/corpus.txt"
	fileKey := fmt.Sprintf("%s/test-list-objects-with-prefix.txt", timestamp)
	_, err := s.UploadFile(ctx, filePath, fileKey, "text/plain", false)
	if err != nil {
		t.Fatal(err)
	}

	objects, err := s.ListObjectsWithPrefix(ctx, timestamp)
	if err != nil {
		t.Fatal(err)
	}

	if len(objects) != 1 {
		t.Fatalf("expected %d objects, got %d", 1, len(objects))
	}

	if objects[0] != fileKey {
		t.Errorf("expected key %q, got %q", fileKey, objects[0])
	}

}

func TestListObjectsWithPrefixStream(t *testing.T) {

	s := New(&Config{
		S3Endpoint:    os.Getenv("YN_S3_ENDPOINT"),
		S3AccessKey:   os.Getenv("YN_S3_ACCESS_KEY"),
		S3SecretKey:   os.Getenv("YN_S3_SECRET_ACCESS_KEY"),
		ContainerName: "test",
	})

	ctx := context.Background()
	filePath := "../../test/corpus.txt"
	timestamp := util.Timestamp(time.Now())
	errs := makeTestErrorStream(t)

	filePathFunc := func(i int) string { return filePath }
	fileKeyFunc := func(i int, s string) string {
		return fmt.Sprintf("%s/test-upload-file-stream-%d.txt", timestamp, i)
	}

	filePathStream := makeTest2StringStream(10, filePathFunc, fileKeyFunc)

	uploadStream := s.UploadFileStream(ctx, errs, filePathStream, false)

	for upload := range uploadStream {
		t.Logf("upload: %s", upload)
	}

	prefixStream := makeTestStringStream(1, func(i int) string {
		return timestamp
	})

	objectStream := s.ListObjectsWithPrefixStream(ctx, errs, prefixStream)

	total := 0
	for object := range objectStream {
		t.Log(object)
		total++
	}

	if total != 10 {
		t.Errorf("expected %d objects, got %d", 10, total)
	}

}
