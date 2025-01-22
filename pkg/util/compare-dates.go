package util

import "time"

func CompareDates(t1, t2 time.Time) int {
	y1, m1, d1 := t1.Date()
	t1 = time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)

	y2, m2, d2 := t2.Date()
	t2 = time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)

	if t1.After(t2) {
		return 1
	}

	if t1.Before(t2) {
		return -1
	}

	return 0
}
