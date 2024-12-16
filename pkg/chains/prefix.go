package chains

import "strings"

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// NewPrefix creates a new prefix from a string containing
// space-separated tokens
func NewPrefix(str string) Prefix {
	return Prefix(strings.Split(str, " "))
}

// First returns the first token in the prefix
func (p Prefix) First() string {
	return p[0]
}

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}
