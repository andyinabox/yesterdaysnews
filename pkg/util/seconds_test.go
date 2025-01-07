package util

import (
	"fmt"
	"testing"
)

func TestSeconds(t *testing.T) {
	expectedStr := "5s"
	seconds := Seconds(5)
	secondsStr := fmt.Sprintf("%s", seconds)

	if secondsStr != "5s" {
		t.Errorf("expected %q, got %q", expectedStr, secondsStr)
	}
}
