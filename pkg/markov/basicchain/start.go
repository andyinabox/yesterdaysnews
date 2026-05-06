package basicchain

import (
	"math/rand"

	"code.andydayton.com/andy/yesterdaysnews/pkg/markov"
	"code.andydayton.com/andy/yesterdaysnews/pkg/markov/basicprefix"
)

func (c *Chain) Start() markov.Prefix {
	return basicprefix.New(c.prefixes[rand.Intn(len(c.prefixes))])
}
