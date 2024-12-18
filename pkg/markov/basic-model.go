package markov

import (
	"encoding/json"
	"io"
	"math/rand"
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

func (m *BasicModel) Build(r io.Reader) error {
	return m.chain.Build(r)
}

func (m *BasicModel) Sentence(p Prefix) (s string, err error) {

	if p == nil {
		p = m.getRandomPrefix()
	}

	words := p.Tokens()

	for i := 0; i < m.maxLength-m.chain.PrefixLength(); i++ {
		next, err := m.chain.Next(p)
		if err != nil {
			break
		}
		words = append(words, next)
	}

	s = strings.Join(words, " ")

	return
}

func (m *BasicModel) Save() ([]byte, error) {
	return json.Marshal(m.chain)
}

func (m *BasicModel) Load(data []byte) error {
	return json.Unmarshal(data, m.chain)
}

func (m *BasicModel) getRandomPrefix() Prefix {
	prefixes := m.chain.Prefixes()
	return NewBasicPrefix(prefixes[rand.Intn(len(prefixes))])
}
