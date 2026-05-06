package basicchain

import (
	"encoding/json"
	"errors"

	"code.andydayton.com/andy/yesterdaysnews/pkg/markov/basicprefix"
)

func (c *Chain) Load(data []byte) error {

	v := struct {
		Version int `json:"version"`
	}{}

	// load based on version. we are ignoring errors here
	// because failure to parse probably means version 0
	json.Unmarshal(data, &v)

	m := Model{}
	switch v.Version {
	case 1:
		err := json.Unmarshal(data, &m)
		if err != nil {
			return err
		}
	// probably the old, unversioned model
	default:
		m.Chain = map[string][]string{}
		err := json.Unmarshal(data, &m.Chain)
		if err != nil {
			return err
		}
	}

	c.version = m.Version
	c.prefixLength = m.PrefixLength
	c.chain = m.Chain

	// re-initializing these to be safe
	c.prefixes = make([]string, len(c.chain))
	c.startPrefixes = []string{}
	c.endChain = make(map[string][]string)

	var i int
	for key, options := range c.chain {

		p := basicprefix.New(key)

		// assume if length is 0, prefix length hasn't been set yet
		if c.prefixLength == 0 {
			c.prefixLength = p.Length()
		}

		// validate prefix length
		if c.prefixLength != p.Length() {
			return errors.New("inconsistent prefix lengths in model")
		}

		// populate start prefixes
		if c.isStartPrefix(key) {
			c.startPrefixes = append(c.startPrefixes, key)
		}

		// build up end tokens
		for _, tok := range options {
			if c.isEndToken(tok) {
				c.endChain[key] = append(c.endChain[key], tok)
			}
		}

		// populate prefixes
		c.prefixes[i] = key
		i++
	}

	return nil
}
