package builder

import (
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/containerservice"
	"gitlab.com/andyinabox/yesterdaysnews/domain/imageprocessor"
	"gitlab.com/andyinabox/yesterdaysnews/domain/videoprocessor"
	"gitlab.com/andyinabox/yesterdaysnews/domain/youtubeservice"
)

type Config struct {
	// env variables
	GoogleAPIKey   string `env:"YN_GOOGLE_API_KEY" required:"true"`
	S3Endpoint     string `env:"YN_S3_ENDPOINT" required:"true"`
	S3AccessKey    string `env:"YN_S3_ACCESS_KEY" required:"true"`
	S3SecretKey    string `env:"YN_S3_SECRET_ACCESS_KEY" required:"true"`
	S3Region       string `env:"YN_S3_REGION" required:"true"`
	BinPathYTDLP   string `env:"YN_YT_DLP_PATH"`
	BinPathFFMPEG  string `env:"YN_FFMPEG_PATH"`
	BinPathFFPROBE string `env:"YN_FFPROBE_PATH"`

	// config variables
	PlaylistIDs                 []string
	ObjectStoreContainerName    string
	OutputDir                   string
	MaxVideoSize                uint
	DownloadCountPerPlaylist    int
	MinClipLengthSeconds        int
	MaxClipLengthSeconds        int
	CaptionPrefixLength         int
	CaptionNewsCorpusWeight     int
	CaptionHospitalCorpusWeight int
	TotalBuildsToKeep           int
	KeepOutputFiles             bool
	SkipUpload                  bool

	// other
	HospitalCorpus string
	OverlayImage   []byte
}

type Builder struct {
	yt   domain.YouTubeService
	vp   domain.VideoProcessor
	ip   domain.ImageProcessor
	cs   domain.ContainerService
	eh   domain.ErrorHandler
	errs chan<- domain.Error
	cfg  *Config
}

func New(cfg *Config, eh domain.ErrorHandler) *Builder {

	yt := youtubeservice.New(&youtubeservice.Config{
		GoogleAPIKey: cfg.GoogleAPIKey,
		BinPathYTDLP: cfg.BinPathYTDLP,
	})

	vp := videoprocessor.New(&videoprocessor.Config{
		FFMpegBinPath:  cfg.BinPathFFMPEG,
		FFProbeBinPath: cfg.BinPathFFPROBE,
	})

	ip := imageprocessor.New()

	cs := containerservice.New(&containerservice.Config{
		S3Endpoint:    cfg.S3Endpoint,
		S3AccessKey:   cfg.S3AccessKey,
		S3SecretKey:   cfg.S3SecretKey,
		ContainerName: cfg.ObjectStoreContainerName,
	})

	return &Builder{
		yt:   yt,
		vp:   vp,
		ip:   ip,
		cs:   cs,
		eh:   eh,
		errs: eh.Channel(),
		cfg:  cfg,
	}
}
