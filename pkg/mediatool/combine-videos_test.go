package mediatool

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCombineVideos(t *testing.T) {
	tool := New("", "")

	t.Run("relative paths", func(t *testing.T) {
		outFile := "../../test/combined-output.mp4"
		t.Cleanup(func() { os.Remove(outFile) })

		result, err := tool.CombineVideos(context.Background(), []string{
			"../../test/video.mp4",
			"../../test/bear.mp4",
		}, outFile)
		if err != nil {
			t.Fatal(err)
		}

		if result != outFile {
			t.Errorf("expected %q, got %q", outFile, result)
		}

		if _, err := os.Stat(outFile); os.IsNotExist(err) {
			t.Error("output file was not created")
		}
	})

	t.Run("absolute paths", func(t *testing.T) {
		video1, err := filepath.Abs("../../test/video.mp4")
		if err != nil {
			t.Fatal(err)
		}
		video2, err := filepath.Abs("../../test/bear.mp4")
		if err != nil {
			t.Fatal(err)
		}
		outFile, err := filepath.Abs("../../test/combined-output-abs.mp4")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { os.Remove(outFile) })

		result, err := tool.CombineVideos(context.Background(), []string{
			video1,
			video2,
		}, outFile)
		if err != nil {
			t.Fatal(err)
		}

		if result != outFile {
			t.Errorf("expected %q, got %q", outFile, result)
		}

		if _, err := os.Stat(outFile); os.IsNotExist(err) {
			t.Error("output file was not created")
		}
	})

	t.Run("empty input", func(t *testing.T) {
		_, err := tool.CombineVideos(context.Background(), []string{}, "output.mp4")
		if err == nil {
			t.Error("expected error for empty input")
		}
	})
}
