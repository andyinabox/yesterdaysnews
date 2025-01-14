package util

import (
	"strconv"
	"time"
)

func Timestamp(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}
