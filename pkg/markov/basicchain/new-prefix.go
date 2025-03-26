package basicchain

import (
	"gitlab.com/andyinabox/yesterdaysnews/pkg/markov"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/markov/basicprefix"
)

func (c *Chain) NewPrefix(s string) markov.Prefix {

	p := basicprefix.New(s)

	if p.Length() != c.PrefixLength() {
		panic(markov.ErrPrefixWrongLength)
	}

	return p
}
