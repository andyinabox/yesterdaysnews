package basicchain

import (
	"math/rand"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/markov"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/markov/basicprefix"
)

func (c *Chain) Start() markov.Prefix {
	return basicprefix.New(c.prefixes[rand.Intn(len(c.prefixes))])
}
