package textprocessor

import (
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/markov"
)

const DefaultPrefixLength = 2

type Config struct {
	MinCaptionLength int
	MaxCaptionLength int
}

type Processor struct {
	chain markov.Chain
	cfg   *Config
}

func New(chain markov.Chain, cfg *Config) *Processor {

	return &Processor{
		chain: chain,
		cfg:   cfg,
	}
}
