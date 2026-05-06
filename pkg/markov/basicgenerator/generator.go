package markov

import (
	"strings"

	"code.andydayton.com/andy/yesterdaysnews/pkg/markov"
)

type Generator struct {
	chain     markov.Chain
	minLength int
	maxLength int
}

func New(chain markov.Chain, minLength int, maxLength int) *Generator {
	return &Generator{chain, minLength, maxLength}
}

func (g *Generator) Sentence(p markov.Prefix) string {

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
