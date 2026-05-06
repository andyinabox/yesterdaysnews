package domain

import "code.andydayton.com/andy/yesterdaysnews/pkg/markov"

type CorpusType string

const (
	CorpusTypeString CorpusType = "string"
	CorpusTypeText   CorpusType = "txt"
	CorpusTypeVTT    CorpusType = "vtt"
)

type Corpus struct {
	Type     CorpusType
	FileGlob string
	Content  string
	Weight   int
}

type CaptionsChain interface {
	markov.Chain
	BuildFromMultiple([]Corpus) error
}
