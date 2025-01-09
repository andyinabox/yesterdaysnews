//go:build objectstoretest
// +build objectstoretest

package objectstoreservice

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}
}

func TestUploadVideo(t *testing.T) {

	s := New(&Config{
		S3Endpoint:    os.Getenv("YN_S3_ENDPOINT"),
		S3AccessKey:   os.Getenv("YN_S3_ACCESS_KEY"),
		S3SecretKey:   os.Getenv("YN_S3_SECRET_ACCESS_KEY"),
		ContainerName: "test",
	})

	ctx := context.Background()

	filePath := "../../test/video.mp4"
	fileKey := fmt.Sprintf("%s.mp4", util.Timestamp(time.Now()))
	file, err := s.UploadFile(ctx, filePath, fileKey, "video/mp4", false)
	if err != nil {
		t.Fatal(err)
	}

	if file != fileKey {
		t.Errorf("expected to get %q for fileKey, got %q", fileKey, file)
	}

}

func TestUploadVideoStream(t *testing.T) {

	filePathStream := make(chan [2]string)
	errs := make(chan domain.Error)
	defer close(errs)

	s := New(&Config{
		S3Endpoint:    os.Getenv("YN_S3_ENDPOINT"),
		S3AccessKey:   os.Getenv("YN_S3_ACCESS_KEY"),
		S3SecretKey:   os.Getenv("YN_S3_SECRET_ACCESS_KEY"),
		ContainerName: "test",
	})

	ctx := context.Background()
	filePath := "../../test/video.mp4"
	fileKey := fmt.Sprintf("%s.mp4", util.Timestamp(time.Now()))

	stream := s.UploadFileStream(ctx, errs, filePathStream, "video/mp4", false)

	// handle errors
	go func() {
		err := <-errs
		t.Error(err)
	}()

	// feed data
	go func() {
		t.Log("start adding to filePathStream")
		filePathStream <- [2]string{filePath, fileKey}
		t.Log("done adding to filePathStream")
		close(filePathStream)
	}()

	for upload := range stream {
		t.Log(upload)
	}

}
