package markov

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrPrefixWrongLength = errors.New("prefix is wrong length")
	ErrNoOptionsFound    = errors.New("no options found for prefix")
)

type Prefix interface {
	fmt.Stringer
	First() string
	Shift(string)
	Length() int
	Tokens() []string
}

type Chain interface {
	Start() Prefix
	Next(Prefix) string
	End(Prefix) string
	Build(r io.Reader)
	PrefixLength() int
	Save() ([]byte, error)
	Load([]byte) error
}

type Model interface {
	Sentence(Prefix) string
	Build(r io.Reader)
	Save() ([]byte, error)
	Load([]byte) error
}
