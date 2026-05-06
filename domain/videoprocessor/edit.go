package videoprocessor

import (
	"time"

	"code.andydayton.com/andy/yesterdaysnews/domain"
	"code.andydayton.com/andy/yesterdaysnews/pkg/mediatool"
)

type edit [2]time.Duration

func newEdit(start, duration mediatool.Duration) domain.VideoEdit {
	return edit{time.Duration(start), time.Duration(duration)}
}

func (e edit) Start() time.Duration {
	return e[0]
}

func (e edit) End() time.Duration {
	return e[0] + e[1]
}

func (e edit) Duration() time.Duration {
	return e[1]
}
