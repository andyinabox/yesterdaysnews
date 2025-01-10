//go:build objectstoretest
// +build objectstoretest

package objectstoreservice

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func TestMoveObject(t *testing.T) {

	s := New(&Config{
		S3Endpoint:    os.Getenv("YN_S3_ENDPOINT"),
		S3AccessKey:   os.Getenv("YN_S3_ACCESS_KEY"),
		S3SecretKey:   os.Getenv("YN_S3_SECRET_ACCESS_KEY"),
		ContainerName: "test",
	})

	ctx := context.Background()
	timestamp := util.Timestamp(time.Now())

	t.Logf("timestamp: %q", timestamp)

	filePath := "../../test/corpus.txt"
	fileKey := fmt.Sprintf("test-move-object-%s-1.txt", timestamp)
	file, err := s.UploadFile(ctx, filePath, fileKey, "video/mp4", false)
	if err != nil {
		t.Fatal(err)
	}

	from := file
	to := fmt.Sprintf("test-move-object-%s-2.txt", timestamp)

	moved, err := s.MoveObject(ctx, from, to)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("moved %q", moved)
}

func TestMoveFileStream(t *testing.T) {

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
		return fmt.Sprintf("test-move-file-stream-%s-%d.txt", timestamp, i)
	}

	filePathStream := makeTest2StringStream(10, filePathFunc, fileKeyFunc)
	uploads := s.UploadFileStream(ctx, errs, filePathStream, "text/plain", false)
	movePathStream := makeTestStringTo2StringStream(uploads, func(s string) string {
		return strings.Replace(s, "test-move-file-stream", "test-move-file-stream-moved", 1)
	})
	movedStream := s.MoveObjectStream(ctx, errs, movePathStream)

	for moved := range movedStream {
		t.Log(moved)
	}

}
