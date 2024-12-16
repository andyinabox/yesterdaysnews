package chains

import (
	"math/rand"
	"strings"

	"github.com/charmbracelet/log"
)

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate() string {
	p := c.getStartPrefix()
	return c.generate(c.cfg.MinWords, c.cfg.MaxWords, p, []string(p))
}

func (c *Chain) GenerateFromToken(tok string) string {

	p, ok := c.getTokenPrefix(tok)

	if !ok {
		p = c.getStartPrefix()
	}

	return c.generate(c.cfg.MinWords, c.cfg.MaxWords, p, []string(p))
}

func (c *Chain) generate(minWords int, maxWords int, startPrefix Prefix, startWords []string) string {
	log.Debug("generate chain", "startPrefix", startPrefix, "startWords", startWords)
	var words []string
	var p Prefix

	outerLoops := 0

	// this helps us keep the word lenth average closer
	// to the middle
	shouldUseEndChoices := func(i int) bool {
		if i < minWords {
			return false
		}

		// the likelyhood of using the end choice will
		// increase as the numbers get largers
		p := float32(i-minWords) / float32(maxWords)
		r := rand.Float32() * 0.9

		return p > r
	}

	// outer loop will keep trying to make a sentence less than n words
outer:
	for {

		outerLoops++

		// reset to start conditions
		words = make([]string, len(startWords))
		copy(words, startWords)
		p = make(Prefix, len(startPrefix))
		copy(p, startPrefix)

		// escape hatch to avoid infinite loops
		if outerLoops > c.cfg.MaxLoops {
			log.Warn("reached max outer loops for generating chain, resetting", "pfx", startPrefix, "startWords", startWords)
			p = c.getStartPrefix()
			startWords = []string{}
		}

		log.Debug("start sentence builing loop", "pfx", p, "words", words)

		// build sentence
		for i := 0; i < maxWords-c.cfg.PrefixLength; i++ {

			// get a sice of possible next tokens
			choices := c.chain[p.String()]
			endChoices := c.endings[p.String()]

			shouldTryForEnd := shouldUseEndChoices(i)

			if len(endChoices) != 0 && shouldTryForEnd {
				choices = endChoices
			}

			// bail out if there are no choices
			if len(choices) == 0 {
				break
			}

			// get a random token
			next := choices[rand.Intn(len(choices))]

			// append that random token to our list of words
			words = append(words, next)

			// finish if the last word has a period
			if shouldTryForEnd && strings.Contains(next, ".") {
				break outer
			}

			// shift token to the end of our prefix before the next loop
			p.Shift(next)
		}

		log.Debug("reached the end of outer loop without sendtence completion. retrying")
	}

	return strings.Join(words, " ")
}
