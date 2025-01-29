package basicchain

import "unicode"

func (c *Chain) isStartPrefix(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsUpper([]rune(s)[0])
}
