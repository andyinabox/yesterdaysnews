package util

import (
	"encoding/base64"
	"time"
)

func TimestampBase64(t time.Time, pad bool) (s string) {

	timestamp := Timestamp(t)

	if pad {
		return base64.URLEncoding.EncodeToString([]byte(timestamp))
	}

	return base64.RawURLEncoding.EncodeToString([]byte(timestamp))
}
