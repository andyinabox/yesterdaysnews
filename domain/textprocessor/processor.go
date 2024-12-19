package textprocessor

import "gitlab.com/andyinabox/yesterdays-news-downloader/pkg/markov"

const DefaultPrefixLength = 2

type Config struct {
	PrefixLength int
}

type Processor struct {
	chain *markov.MultiChain
	cfg   *Config
}

func New(cfg *Config) *Processor {

	if cfg.PrefixLength == 0 {
		cfg.PrefixLength = DefaultPrefixLength
	}

	return &Processor{
		chain: markov.NewMultiChain(cfg.PrefixLength),
		cfg:   cfg,
	}
}
