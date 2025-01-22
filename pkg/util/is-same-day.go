package util

import "time"

func IsSameDay(t1, t2 time.Time) bool {
	return CompareDates(t1, t2) == 0
}
