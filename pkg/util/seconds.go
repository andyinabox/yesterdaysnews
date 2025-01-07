package util

import "time"

func Seconds(seconds int) time.Duration {
	return time.Duration(int(time.Second) * seconds)
}
