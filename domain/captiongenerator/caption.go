package captiongenerator

import (
	"strings"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/markov"
)

func (g *Generator) Caption(prevSentence string) string {

	prefix := g.getPrefixFromSentence(prevSentence)

	words := []string{}

	prefixLen := g.chain.PrefixLength()
	minLen := g.cfg.MinCaptionLength
	maxLen := g.cfg.MaxCaptionLength

	for i := 0; i < maxLen-prefixLen; i++ {
		var next string

		// try for an ending
		if i > minLen-prefixLen {
			next = g.chain.End(prefix)
			if next != "" {
				words = append(words, next)
				break
			}
		}

		// try for a regular word
		next = g.chain.Next(prefix)
		if next == "" {
			break
		}
		words = append(words, next)

		prefix.Shift(next)
	}

	return strings.Join(words, " ")
}

func (g *Generator) getPrefixFromSentence(s string) markov.Prefix {
	if s == "" {
		return g.chain.Start()
	}

	tokens := strings.Split(s, " ")
	if len(tokens) < g.chain.PrefixLength() {
		return g.chain.Start()
	}

	str := tokens[len(tokens)-g.chain.PrefixLength():]

	return g.chain.NewPrefix(strings.Join(str, " "))
}
