package basicchain

import "encoding/json"

func (c *Chain) Save() ([]byte, error) {
	return json.Marshal(c.model)
}
