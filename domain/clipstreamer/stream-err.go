package clipstreamer

import "gitlab.com/andyinabox/yesterdaysnews/domain"

type StreamErr struct {
	typ domain.StreamErrType
	err error
}

func NewStreamErr(t domain.StreamErrType, err error) domain.StreamErr {
	return &StreamErr{t, err}
}

func (e *StreamErr) Error() string {
	return e.err.Error()
}

func (e *StreamErr) Type() domain.StreamErrType {
	return e.typ
}
