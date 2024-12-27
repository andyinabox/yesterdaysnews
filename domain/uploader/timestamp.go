package uploader

import (
	"strconv"
	"time"
)

func timestampFromTime(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

func timestamp() string {
	return timestampFromTime(time.Now())
}
