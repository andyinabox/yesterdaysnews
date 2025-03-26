package youtubedownloader

import (
	"context"
	"path/filepath"
	"testing"
)

func TestCheckVersion(t *testing.T) {

	{
		binPath, err := filepath.Abs("../../test/bin/yt-dlp-equalversion")
		if err != nil {
			t.Fatal(err)
		}
		d := New(binPath)
		err = d.CheckVersion(context.Background())
		if err != nil {
			t.Error(err)
		}
	}
	{
		binPath, err := filepath.Abs("../../test/bin/yt-dlp-higherversion")
		if err != nil {
			t.Fatal(err)
		}
		d := New(binPath)
		err = d.CheckVersion(context.Background())
		if err != nil {
			t.Error(err)
		}
	}
	{
		binPath, err := filepath.Abs("../../test/bin/yt-dlp-incompatversion")
		if err != nil {
			t.Fatal(err)
		}
		d := New(binPath)
		err = d.CheckVersion(context.Background())
		if err == nil {
			t.Error("error expected")
		}
	}
}
