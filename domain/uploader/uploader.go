package uploader

type Config struct {
	VideoFileName    string
	SubsFileName     string
	CombinedFileName string
	ModelFileName    string
	ClipsDirName     string
}

type Uploader struct {
	cfg *Config
}

func New(cfg *Config) *Uploader {
	return &Uploader{cfg}
}
