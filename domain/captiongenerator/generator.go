package captiongenerator

import (
	"gitlab.com/andyinabox/yesterdaysnews/pkg/markov"
)

const DefaultPrefixLength = 2

type Config struct {
	MinCaptionLength int
	MaxCaptionLength int
}

type Generator struct {
	chain markov.Chain
	cfg   *Config
}

func New(chain markov.Chain, cfg *Config) *Generator {

	return &Generator{
		chain: chain,
		cfg:   cfg,
	}
}

func (g *Generator) MinCaptionLength() int {
	return g.cfg.MinCaptionLength
}

func (g *Generator) MaxCaptionLength() int {
	return g.cfg.MaxCaptionLength
}
