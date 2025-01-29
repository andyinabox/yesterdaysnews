package basicchain

import (
	"bufio"
	"fmt"
	"io"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/markov/basicprefix"
)

func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(basicprefix.Prefix, c.prefixLength)

	for {
		var s string

		// scan next word into string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}

		// get prefix key
		key := p.String()

		// add to main chain
		c.chain[key] = append(c.chain[key], s)

		// populate endings chain
		if c.isEndToken(s) {
			c.endChain[key] = append(c.endChain[key], s)
		}

		// shift the prefix
		p.Shift(s)
	}

	c.prefixes = make([]string, len(c.chain))
	var i int
	for k := range c.chain {

		// populate start prefixes
		if c.isStartPrefix(k) {
			c.startPrefixes = append(c.startPrefixes, k)
		}

		// populate prefixes
		c.prefixes[i] = k
		i++
	}
}
