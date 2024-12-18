package textprocessor

type Config struct{}

type Processor struct {
	cfg *Config
}

func New(cfg *Config) *Processor {
	return &Processor{cfg}
}
