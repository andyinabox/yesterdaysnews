package youtubedownloader

import (
	"context"
	"testing"
)

func TestGetVideoInfo(t *testing.T) {
	// t.Skip()

	id := "8N5yiQ1SABE"

	client := New("")

	video, err := client.GetVideoInfo(context.Background(), id, "")
	if err != nil {
		t.Fatal(err)
	}

	if video.ID != id {
		t.Errorf("expected id %q, got %q", id, video.ID)
	}
}
