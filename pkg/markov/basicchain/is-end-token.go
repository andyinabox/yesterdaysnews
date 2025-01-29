package basicchain

import "strings"

func (c *Chain) isEndToken(s string) bool {
	return strings.HasSuffix(s, ".")
}
