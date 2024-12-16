package chains

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.cfg.PrefixLength)

	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}

		key := p.String()

		// add to main chain
		c.chain[key] = append(c.chain[key], s)

		// add to endings chain if ends in a period
		if strings.HasSuffix(s, ".") {
			c.endings[key] = append(c.endings[key], s)
		}

		// add to start prefixes if first char is capital letter
		if unicode.IsUpper([]rune(key)[0]) {
			c.startPrefixes = append(c.startPrefixes, key)
		}

		// add key to a lookup which will allow looking up a prefix
		// key using the first token in the prefix
		first := p.First()
		c.prefixLookup[first] = append(c.prefixLookup[first], key)

		p.Shift(s)
	}
}

func (c *Chain) BuildFromString(s string) {
	c.Build(strings.NewReader(s))
}
