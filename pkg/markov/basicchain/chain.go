package basicchain

const modelVersion = 1

type Model struct {
	Version      int                 `json:"version"`
	Chain        map[string][]string `json:"chain"`
	PrefixLength int                 `json:"prefixLength"`
}

type Chain struct {
	// model version
	version int
	// main chain, keys are prefixes, values are tokens
	chain map[string][]string
	// number of tokens in a prefix
	prefixLength int
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
		version:       modelVersion,
		chain:         make(map[string][]string),
		prefixLength:  prefixLength,
		prefixes:      []string{},
		startPrefixes: []string{},
		endChain:      make(map[string][]string),
	}
}
