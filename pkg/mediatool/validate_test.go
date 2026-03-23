package mediatool

import (
	"context"
	"testing"
)

func TestValidate(t *testing.T) {
	ctx := context.Background()
	mt := New("", "")
	// test valid video
	{
		path := "../../test/encodingerr/clips/goodvideo.webm"
		err := mt.Validate(ctx, path)
		if err != nil {
			t.Error(err)
		}
	}

	// test invalid video
	{
		path := "../../test/encodingerr/clips/badvideo.webm"
		err := mt.Validate(ctx, path)
		t.Log(err)
		if err == nil {
			t.Error("expected error, got none")
		}
	}
}
