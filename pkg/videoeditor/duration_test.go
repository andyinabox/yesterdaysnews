package videoeditor

import (
	"testing"
	"time"
)

func TestDurationToTimestamp(t *testing.T) {
	tries := map[string]string{
		"1h34m": "01:34:00",
		"25s":   "00:00:25",
		"1m12s": "00:01:12",
	}

	for str, expected := range tries {
		duration, err := time.ParseDuration(str)
		if err != nil {
			t.Fatal(err)
		}
		d := Duration(duration)
		result := d.String()
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	}
}
