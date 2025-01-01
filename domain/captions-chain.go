package domain

import "gitlab.com/andyinabox/yesterdaysnews/pkg/markov"

type CorpusType string

const (
	CorpusTypeText CorpusType = "txt"
	CorpusTypeVTT  CorpusType = "vtt"
)

type Corpus struct {
	Type     CorpusType
	FileGlob string
	Weight   int
}

type CaptionsChain interface {
	markov.Chain
	BuildFromMultiple([]Corpus) error
}
