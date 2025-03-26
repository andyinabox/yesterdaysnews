package basicchain

import (
	"math/rand"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/markov"
)

func (c *Chain) End(p markov.Prefix) string {
	// validate length
	if p.Length() != c.prefixLength {
		panic(markov.ErrPrefixWrongLength)
	}

	options, ok := c.endChain[p.String()]
	if !ok {
		return ""
	}

	return options[rand.Intn(len(options))]
}
