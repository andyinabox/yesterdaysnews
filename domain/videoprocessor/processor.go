package videoprocessor

import "gitlab.com/andyinabox/yesterdays-news-downloader/pkg/mediatool"

type Config struct {
	FFMpegBinPath  string
	FFProbeBinPath string
}

type Processor struct {
	mt  *mediatool.Tool
	cfg *Config
}

func New(cfg *Config) *Processor {
	return &Processor{
		mt: mediatool.New(cfg.FFMpegBinPath, cfg.FFProbeBinPath),
	}
}
