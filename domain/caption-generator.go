package domain

type CaptionGenerator interface {
	Caption(string) string
	MinCaptionLength() int
	MaxCaptionLength() int
}
