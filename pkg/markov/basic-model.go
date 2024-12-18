package markov

import (
	"encoding/json"
	"io"
	"strings"
)

type BasicModel struct {
	chain     Chain
	minLength int
	maxLength int
}

func NewBasicModel(prefixLength int, minLength int, maxLength int) *BasicModel {
	return &BasicModel{
		chain: NewBasicChain(prefixLength),
	}
}

func (m *BasicModel) Build(r io.Reader) {
	m.chain.Build(r)
}

func (m *BasicModel) Sentence(p Prefix) string {

	if p == nil {
		p = m.chain.Start()
	}

	words := p.Tokens()

	for i := 0; i < m.maxLength-m.chain.PrefixLength(); i++ {
		var next string

		// try for an ending
		if i > m.minLength-m.chain.PrefixLength() {
			next = m.chain.End(p)
			if next != "" {
				words = append(words, next)
				break
			}
		}

		// try for a regular word
		next = m.chain.Next(p)
		if next == "" {
			break
		}
		words = append(words, next)
	}

	return strings.Join(words, " ")
}

func (m *BasicModel) Save() ([]byte, error) {
	return json.Marshal(m.chain)
}

func (m *BasicModel) Load(data []byte) error {
	return json.Unmarshal(data, m.chain)
}
