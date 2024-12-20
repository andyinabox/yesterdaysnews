package markov

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"unicode"

	"github.com/charmbracelet/log"
)

type BasicChain struct {
	// main chain, keys are prefixes, values are tokens
	chain map[string][]string
	// prefixes is just a duplicate of the chain keys
	prefixes []string
	// startPrefixes are prefixes that can be used to start a sentence
	startPrefixes []string
	// a subset of the chain that only includes ending words
	endChain map[string][]string

	prefixLength int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewBasicChain(prefixLength int) *BasicChain {
	return &BasicChain{
		chain:         make(map[string][]string),
		prefixes:      []string{},
		startPrefixes: []string{},
		endChain:      make(map[string][]string),

		prefixLength: prefixLength,
	}
}

func (c *BasicChain) Next(p Prefix) string {

	// validate length
	if p.Length() != c.prefixLength {
		log.Fatal(ErrPrefixWrongLength)
	}

	options, ok := c.chain[p.String()]
	if !ok {
		return ""
	}
	return options[rand.Intn(len(options))]
}

func (c *BasicChain) PrefixLength() int {
	return c.prefixLength
}

func (c *BasicChain) NewPrefix(s string) Prefix {

	p := NewBasicPrefix(s)

	if p.Length() != c.PrefixLength() {
		log.Fatal(ErrPrefixWrongLength)
	}

	return p
}

func (c *BasicChain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(BasicPrefix, c.prefixLength)

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

func (c *BasicChain) Start() Prefix {
	return NewBasicPrefix(c.prefixes[rand.Intn(len(c.prefixes))])
}

func (c *BasicChain) End(p Prefix) string {
	// validate length
	if p.Length() != c.prefixLength {
		log.Fatal(ErrPrefixWrongLength)
	}

	options, ok := c.endChain[p.String()]
	if !ok {
		return ""
	}

	return options[rand.Intn(len(options))]
}

func (c *BasicChain) Save() ([]byte, error) {
	return json.Marshal(c.chain)
}

func (c *BasicChain) Load(data []byte) error {

	// import chain
	err := json.Unmarshal(data, &c.chain)
	if err != nil {
		return err
	}

	// re-initializing these to be safe
	c.prefixes = make([]string, len(c.chain))
	c.startPrefixes = []string{}
	c.endChain = make(map[string][]string)

	var i int
	for key, options := range c.chain {

		p := NewBasicPrefix(key)

		// assume if length is 0, prefix length hasn't been set yet
		if c.prefixLength == 0 {
			c.prefixLength = p.Length()
		}

		// validate prefix length
		if c.prefixLength != p.Length() {
			return errors.New("inconsistent prefix lengths in model")
		}

		// populate start prefixes
		if c.isStartPrefix(key) {
			c.startPrefixes = append(c.startPrefixes, key)
		}

		// build up end tokens
		for _, tok := range options {
			if c.isEndToken(tok) {
				c.endChain[key] = append(c.endChain[key], tok)
			}
		}

		// populate prefixes
		c.prefixes[i] = key
		i++
	}

	return nil
}

func (c *BasicChain) isStartPrefix(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsUpper([]rune(s)[0])
}

func (c *BasicChain) isEndToken(s string) bool {
	return strings.HasSuffix(s, ".")
}
