package videoprocessor

import (
	"context"
	"testing"
	"time"
)

func TestGetVideoEditPoints(t *testing.T) {

	vp := New(&Config{})

	video := "../../test/video.mp4"

	edits, err := vp.GetVideoEditPoints(context.Background(), video, 1*time.Second, 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if len(edits) == 0 {
		t.Error("expected multiple edit points, got zero")
	}

	for _, e := range edits {

		if int(e.Start()) < 0 || int(e.Duration()) < 0 {
			t.Errorf("recieved negative value: %v", e)
			continue
		}
		t.Logf("%v, %v\n", e.Start(), e.End())
	}

}
