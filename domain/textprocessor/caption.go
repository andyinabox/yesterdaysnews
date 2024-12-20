package textprocessor

import (
	"strings"

	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/markov"
)

func (p *Processor) Caption(prevSentence string) string {

	prefix := p.getPrefixFromSentence(prevSentence)

	words := []string{}

	prefixLen := p.chain.PrefixLength()
	minLen := p.cfg.MinCaptionLength
	maxLen := p.cfg.MaxCaptionLength

	for i := 0; i < maxLen-prefixLen; i++ {
		var next string

		// try for an ending
		if i > minLen-prefixLen {
			next = p.chain.End(prefix)
			if next != "" {
				words = append(words, next)
				break
			}
		}

		// try for a regular word
		next = p.chain.Next(prefix)
		if next == "" {
			break
		}
		words = append(words, next)

		prefix.Shift(next)
	}

	return strings.Join(words, " ")
}

func (p *Processor) getPrefixFromSentence(s string) markov.Prefix {
	if s == "" {
		return p.chain.Start()
	}

	tokens := strings.Split(s, " ")
	if len(tokens) < p.chain.PrefixLength() {
		return p.chain.Start()
	}

	str := tokens[len(tokens)-p.chain.PrefixLength():]

	return p.chain.NewPrefix(strings.Join(str, " "))
}
