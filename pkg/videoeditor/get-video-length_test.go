package videoeditor

import (
	"context"
	"testing"
)

func TestGetVideoLength(t *testing.T) {

	path := "../../test/video.mp4"
	expected := "00:00:11"

	// TODO: find out why godotenv isn't working here
	e := New("", "")

	result, err := e.GetVideoLength(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}

	if expected != result.String() {
		t.Errorf("expected %q, got %q", expected, result.String())
	}
}
