package videoeditor

import (
	"fmt"
	"time"
)

type Duration time.Duration

func (d Duration) Timestamp() string {
	dur := time.Duration(d).Round(time.Second)
	h := dur / time.Hour
	dur -= h * time.Hour
	m := dur / time.Minute
	dur -= m * time.Minute
	s := dur / time.Second

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
