package basicprefix

import (
	"strings"
)

type Prefix []string

func New(str string) Prefix {
	return Prefix(strings.Split(str, " "))
}

func (p Prefix) First() string {
	return p[0]
}

func (p Prefix) String() string {
	return strings.Join(p, " ")
}

func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

func (p Prefix) Length() int {
	return len(p)
}

func (p Prefix) Tokens() []string {
	return []string(p)
}
