//go:build objectstoretest
// +build objectstoretest

package containerservice

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}
}

func TestUploadFile(t *testing.T) {

	s := New(&Config{
		S3Endpoint:    os.Getenv("YN_S3_ENDPOINT"),
		S3AccessKey:   os.Getenv("YN_S3_ACCESS_KEY"),
		S3SecretKey:   os.Getenv("YN_S3_SECRET_ACCESS_KEY"),
		S3Region:      os.Getenv("YN_S3_REGION"),
		ContainerName: "yesterdaysnewstest",
	})

	ctx := context.Background()

	filePath := "../../test/corpus.txt"
	fileKey := fmt.Sprintf("test-upload-file-%s.txt", util.Timestamp(time.Now()))
	file, err := s.UploadFile(ctx, filePath, fileKey, "text/plain", false)
	if err != nil {
		t.Fatal(err)
	}

	if file != fileKey {
		t.Errorf("expected to get %q for fileKey, got %q", fileKey, file)
	}

}

func TestUploadFileStream(t *testing.T) {

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
		return fmt.Sprintf("test-upload-file-stream-%s-%d.txt", timestamp, i)
	}

	filePathStream := makeTest2StringStream(10, filePathFunc, fileKeyFunc)

	stream := s.UploadFileStream(ctx, errs, filePathStream, "text/plain", false)

	for upload := range stream {
		t.Log(upload)
	}

}
