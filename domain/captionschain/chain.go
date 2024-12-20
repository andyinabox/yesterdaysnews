package captionschain

import "gitlab.com/andyinabox/yesterdays-news-downloader/pkg/markov"

type Chain struct {
	markov.BasicChain
}

func New(prefixLength int) *Chain {
	return &Chain{
		BasicChain: *markov.NewBasicChain(prefixLength),
	}
}
