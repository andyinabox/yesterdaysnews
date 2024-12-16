package videoeditor

import (
	"context"
	"testing"
)

func init() {
	// godotenv.Load()
	// log.SetLevel(log.DebugLevel)
}

func TestGetVideoLength(t *testing.T) {

	path := "../../test/video.mp4"
	expected := "11.011s"

	// TODO: find out why godotenv isn't working here
	e := New("/opt/homebrew/bin/ffmpeg", "/opt/homebrew/bin/ffprobe")

	result, err := e.GetVideoLength(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}

	if expected != result.String() {
		t.Errorf("expected %q, got %q", expected, result.String())
	}
}
