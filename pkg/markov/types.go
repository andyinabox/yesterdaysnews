package markov

import (
	"encoding/json"
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
	json.Marshaler
	json.Unmarshaler
	Prefixes() []string
	Next(Prefix) (string, error)
	Build(r io.Reader) error
	PrefixLength() int
}

type Model interface {
	Sentence(Prefix) (string, error)
	Build(r io.Reader) error
	Save() ([]byte, error)
	Load([]byte) error
}
