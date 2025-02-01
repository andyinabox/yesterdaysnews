package util

import (
	"testing"
	"time"
)

func TestTimestampBase64(t *testing.T) {
	date := time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC)
	expectedPad := "MTcwNDA2NzIwMA=="
	expectedNoPad := "MTcwNDA2NzIwMA"

	ts := TimestampBase64(date, true)
	t.Log(ts)

	if expectedPad != ts {
		t.Errorf("expected %q, got %q", expectedPad, ts)
	}

	ts = TimestampBase64(date, false)
	t.Log(ts)

	if expectedNoPad != ts {
		t.Errorf("expected %q, got %q", expectedNoPad, ts)
	}

}
