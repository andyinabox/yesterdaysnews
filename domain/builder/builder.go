package builder

import (
	"time"

	"code.andydayton.com/andy/yesterdaysnews/domain"
	"code.andydayton.com/andy/yesterdaysnews/domain/containerservice"
	"code.andydayton.com/andy/yesterdaysnews/domain/imageprocessor"
	"code.andydayton.com/andy/yesterdaysnews/domain/videoprocessor"
	"code.andydayton.com/andy/yesterdaysnews/domain/youtubeservice"
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
	BuildID                     string
	PlaylistIDs                 []string
	ObjectStoreContainerName    string
	OutputDir                   string
	MaxVideoSize                uint
	DownloadCountPerPlaylist    int
	MaxPlaylistRequests         int
	MinClipLengthSeconds        int
	MaxClipLengthSeconds        int
	CaptionPrefixLength         int
	CaptionNewsCorpusWeight     int
	CaptionHospitalCorpusWeight int
	CaptionMinDuration          float64
	CaptionMaxDuration          float64
	TotalBuildsToKeep           int
	KeepOutputFiles             bool
	SkipUpload                  bool
	ThrottleDownloadsBy         time.Duration

	// other
	HospitalCorpus string
	OverlayImage   []byte
}

type Builder struct {
	yt   domain.YouTubeService
	vp   domain.VideoProcessor
	ip   domain.ImageProcessor
	cs   domain.ContainerService
	cg   domain.CaptionGenerator
	eh   domain.ErrorHandler
	errs chan<- domain.Error
	cfg  *Config
}

func New(cfg *Config, eh domain.ErrorHandler) domain.Builder {

	yt := youtubeservice.New(&youtubeservice.Config{
		GoogleAPIKey:        cfg.GoogleAPIKey,
		BinPathYTDLP:        cfg.BinPathYTDLP,
		MaxPlaylistRequests: cfg.MaxPlaylistRequests,
		ThrottleDownloadsBy: cfg.ThrottleDownloadsBy,
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
