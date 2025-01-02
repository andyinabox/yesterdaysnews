package clipstreamer

import "gitlab.com/andyinabox/yesterdaysnews/domain"

type Err struct {
	typ domain.StreamErrType
	err error
}

func NewErr(t domain.StreamErrType, err error) domain.StreamErr {
	return &Err{t, err}
}

func (e *Err) Error() string {
	return e.err.Error()
}

func (e *Err) Type() domain.StreamErrType {
	return e.typ
}
