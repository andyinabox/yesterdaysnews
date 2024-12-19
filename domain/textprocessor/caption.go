package textprocessor

import "strings"

func (p *Processor) Caption() string {

	prefix := p.chain.Start()

	words := prefix.Tokens()

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
