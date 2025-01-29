package basicchain

import (
	"math/rand"
	"strings"
	"unicode"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/markov"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/markov/basicprefix"
)

const modelVersion = 1

type model struct {
	// model version
	version int
	// main chain, keys are prefixes, values are tokens
	chain map[string][]string
	// number of tokens in a prefix
	prefixLength int
}

type Chain struct {
	model
	// prefixes is just a duplicate of the chain keys
	prefixes []string
	// startPrefixes are prefixes that can be used to start a sentence
	startPrefixes []string
	// a subset of the chain that only includes ending words
	endChain map[string][]string
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func New(prefixLength int) *Chain {
	return &Chain{
		model: model{
			version:      modelVersion,
			chain:        make(map[string][]string),
			prefixLength: prefixLength,
		},
		prefixes:      []string{},
		startPrefixes: []string{},
		endChain:      make(map[string][]string),
	}
}

func (c *Chain) Next(p markov.Prefix) string {

	// validate length
	if p.Length() != c.prefixLength {
		log.Fatal(markov.ErrPrefixWrongLength)
	}

	options, ok := c.chain[p.String()]
	if !ok {
		return ""
	}
	return options[rand.Intn(len(options))]
}

func (c *Chain) PrefixLength() int {
	return c.prefixLength
}

func (c *Chain) NewPrefix(s string) markov.Prefix {

	p := basicprefix.New(s)

	if p.Length() != c.PrefixLength() {
		log.Fatal(markov.ErrPrefixWrongLength)
	}

	return p
}

func (c *Chain) Start() markov.Prefix {
	return basicprefix.New(c.prefixes[rand.Intn(len(c.prefixes))])
}

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

func (c *Chain) isStartPrefix(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsUpper([]rune(s)[0])
}

func (c *Chain) isEndToken(s string) bool {
	return strings.HasSuffix(s, ".")
}
