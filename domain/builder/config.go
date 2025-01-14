package builder

type Config struct {
	// env variables
	GoogleAPIKey string `env:"GOOGLE_API_KEY" required:"true"`
	S3Endpoint   string `env:"YN_S3_ENDPOINT" required:"true"`
	S3AccessKey  string `env:"YN_S3_ACCESS_KEY" required:"true"`
	S3SecretKey  string `env:"YN_S3_SECRET_ACCESS_KEY" required:"true"`

	// config variables
	PlaylistIDs              []string
	DownloadCountPerPlaylist int
	ObjectStoreContainerName string
	ObjectStorePrimaryDir    string
	OutputDir                string
	ClipsDirName             string
	MinClipLengthSeconds     int
	MaxClipLengthSeconds     int
	CaptionPrefixLength      int
	RemoveFilesOnCompletion  bool
	TotalPrefixesToKeep      int
}
