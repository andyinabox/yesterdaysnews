package markov

import (
	"strings"
)

type BasicGenerator struct {
	chain     Chain
	minLength int
	maxLength int
}

func NewBasicGenerator(chain Chain, minLength int, maxLength int) *BasicGenerator {
	return &BasicGenerator{chain, minLength, maxLength}
}

func (g *BasicGenerator) Sentence(p Prefix) string {

	if p == nil {
		p = g.chain.Start()
	}

	words := p.Tokens()

	for i := 0; i < g.maxLength-g.chain.PrefixLength(); i++ {
		var next string

		// try for an ending
		if i > g.minLength-g.chain.PrefixLength() {
			next = g.chain.End(p)
			if next != "" {
				words = append(words, next)
				break
			}
		}

		// try for a regular word
		next = g.chain.Next(p)
		if next == "" {
			break
		}
		words = append(words, next)

		p.Shift(next)
	}

	return strings.Join(words, " ")
}
