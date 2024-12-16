package chains

import (
	"math/rand"

	"github.com/charmbracelet/log"
)

const DefaultMaxLoops = 10

type Config struct {
	PrefixLength int
	MinWords     int
	MaxWords     int
	MaxLoops     int
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	// main markov chain. note that the key is actuall a string representation of a prefix
	chain map[string][]string
	// chain of prefixes that lead to end tokens
	endings map[string][]string
	// contains prefixes that begin with a capital latter (stored as strings)
	startPrefixes []string
	// the key here is a single token, used to lookup possible prefixes that begin
	// with the given token
	prefixLookup map[string][]string
	cfg          *Config
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(cfg *Config) *Chain {

	if cfg.MaxLoops == 0 {
		cfg.MaxLoops = DefaultMaxLoops
	}

	log.Debug("create new chain", "config", cfg)

	return &Chain{
		chain:         make(map[string][]string),
		endings:       make(map[string][]string),
		startPrefixes: make([]string, 0),
		prefixLookup:  make(map[string][]string),
		cfg:           cfg,
	}
}

func (c *Chain) getStartPrefix() Prefix {
	p1 := make(Prefix, c.cfg.PrefixLength)
	p2 := NewPrefix(c.startPrefixes[rand.Intn(len(c.startPrefixes))])
	copy(p1, p2)
	return p1
}

func (c *Chain) getTokenPrefix(tok string) (p1 Prefix, ok bool) {

	prefixes, ok := c.prefixLookup[tok]

	if !ok || len(prefixes) < 1 {
		return p1, false
	}

	// copy prefix
	p1 = make(Prefix, c.cfg.PrefixLength)
	p2 := NewPrefix(prefixes[rand.Intn(len(prefixes))])
	copy(p1, p2)
	ok = true

	return
}
