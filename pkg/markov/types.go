package markov

import (
	"errors"
	"fmt"
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
	PrefixLength() int
	Save() ([]byte, error)
	Load([]byte) error
}

type Generator interface {
	Sentence(Prefix) string
}
