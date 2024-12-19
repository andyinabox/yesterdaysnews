package textprocessor

import (
	"errors"

	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/markov"
)

func (p *Processor) Sentence(pre markov.Prefix) (string, error) {
	return "", errors.New("not implemented")
}
