package errorhandler

import (
	"fmt"
)

type Error interface {
	error
	Type() string
}

type err struct {
	typ string
	err error
}

func (e *err) Type() string {
	return e.typ
}

func (e *err) Error() string {
	return e.err.Error()
}

func (e *err) String() string {
	return fmt.Sprintf("%s: %s", e.typ, e.err.Error())
}

func Err(typ string, e error) Error {
	return &err{typ, e}
}
