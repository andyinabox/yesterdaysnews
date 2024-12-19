package markov

import (
	"strings"
)

type BasicPrefix []string

func NewBasicPrefix(str string) BasicPrefix {
	return BasicPrefix(strings.Split(str, " "))
}

func (p BasicPrefix) First() string {
	return p[0]
}

func (p BasicPrefix) String() string {
	return strings.Join(p, " ")
}

func (p BasicPrefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

func (p BasicPrefix) Length() int {
	return len(p)
}

func (p BasicPrefix) Tokens() []string {
	return []string(p)
}
