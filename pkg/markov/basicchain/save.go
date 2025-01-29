package basicchain

import (
	"encoding/json"
)

func (c *Chain) Save() ([]byte, error) {
	model := Model{
		Version:      c.version,
		PrefixLength: c.prefixLength,
		Chain:        c.chain,
	}

	return json.Marshal(model)
}
