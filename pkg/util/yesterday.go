package util

import "time"

func Yesterday() time.Time {
	return time.Now().AddDate(0, 0, -1)
}
