package util

import "time"

func DateString(t time.Time) string {
	return t.Format("2006-01-02")
}
