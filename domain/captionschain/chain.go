package captionschain

import (
	"code.andydayton.com/andy/yesterdaysnews/pkg/markov/basicchain"
)

type Chain struct {
	basicchain.Chain
}

func New(prefixLength int) *Chain {
	return &Chain{
		Chain: *basicchain.New(prefixLength),
	}
}
