package basicchain

import (
	"math/rand"

	"code.andydayton.com/andy/yesterdaysnews/pkg/markov"
)

func (c *Chain) Next(p markov.Prefix) string {

	// validate length
	if p.Length() != c.prefixLength {
		panic(markov.ErrPrefixWrongLength)
	}

	options, ok := c.chain[p.String()]
	if !ok {
		return ""
	}
	return options[rand.Intn(len(options))]
}
