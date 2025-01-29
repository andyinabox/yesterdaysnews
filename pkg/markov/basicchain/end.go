package basicchain

import (
	"math/rand"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/markov"
)

func (c *Chain) End(p markov.Prefix) string {
	// validate length
	if p.Length() != c.prefixLength {
		log.Fatal(markov.ErrPrefixWrongLength)
	}

	options, ok := c.endChain[p.String()]
	if !ok {
		return ""
	}

	return options[rand.Intn(len(options))]
}
