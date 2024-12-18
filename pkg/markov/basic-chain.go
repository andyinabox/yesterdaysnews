package markov

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
)

type BasicChain struct {
	chain        map[string][]string
	prefixLength int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewBasicChain(prefixLength int) *BasicChain {
	return &BasicChain{
		chain:        make(map[string][]string),
		prefixLength: prefixLength,
	}
}

func (c *BasicChain) Next(p Prefix) (s string, err error) {

	// validate length
	if p.Length() != c.prefixLength {
		err = ErrPrefixWrongLength
		return
	}

	options, ok := c.chain[p.String()]
	if !ok {
		err = ErrNoOptionsFound
		return
	}

	s = options[rand.Intn(len(options))]
	return
}

func (c *BasicChain) PrefixLength() int {
	return c.prefixLength
}

func (c *BasicChain) Build(r io.Reader) error {
	br := bufio.NewReader(r)
	p := make(BasicPrefix, c.prefixLength)

	for {
		var s string

		// scan next word into string
		_, err := fmt.Fscan(br, &s)
		if err != nil {
			return err
		}

		// get prefix key
		key := p.String()

		// add to main chain
		c.chain[key] = append(c.chain[key], s)

		// shift the prefix
		p.Shift(s)
	}

}

func (c *BasicChain) Prefixes() []string {
	prefixes := make([]string, len(c.chain))
	var i int
	for k := range c.chain {
		prefixes[i] = k
		i++
	}
	return prefixes
}

func (c *BasicChain) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.chain)
}

func (c *BasicChain) UnmarshalJSON(data []byte) error {

	err := json.Unmarshal(data, c.chain)
	if err != nil {
		return err
	}

	for k := range c.chain {
		p := NewBasicPrefix(k)

		// assume if length is 0, prefix length hasn't been set yet
		if c.prefixLength == 0 {
			c.prefixLength = p.Length()
			continue
		}

		if c.prefixLength != p.Length() {
			return errors.New("inconsistent prefix lengths in model")
		}
	}

	return nil
}
