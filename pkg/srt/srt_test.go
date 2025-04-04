package srt

import (
	"testing"
	"time"
)

func TestSRT(t *testing.T) {
	expected := `0
00:00:00,000 --> 00:00:05,000
This is the first caption

1
00:00:05,000 --> 00:00:10,000
This is the second caption

`

	s := New()

	dur, err := time.ParseDuration("5s")
	if err != nil {
		t.Fatal(err)
	}
	s.AddToEnd(dur, "This is the first caption")
	s.AddToEnd(dur, "This is the second caption")

	result := s.String()

	if result != expected {
		t.Log(expected)
		t.Log(result)
		t.Errorf("unexpected result")
	}

}
