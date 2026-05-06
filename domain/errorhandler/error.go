package errorhandler

import (
	"fmt"

	"code.andydayton.com/andy/yesterdaysnews/domain"
)

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

func (e *err) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", e.Error())), nil
}

func Err(typ string, e error) domain.Error {
	return &err{typ, e}
}
