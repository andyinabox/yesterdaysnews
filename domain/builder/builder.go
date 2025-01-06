package builder

type Config struct {
	// env variables
	GoogleAPIKey string `env:"GOOGLE_API_KEY"`
	S3Endpoint   string `env:"YN_S3_ENDPOINT"`
	S3AccessKey  string `env:"YN_S3_ACCESS_KEY"`
	S3SecretKey  string `env:"YN_S3_SECRET_ACCESS_KEY"`

	// config variables
	ObjectStoreContainerName string
	ObjectStorePrimaryDir    string
	DownloadDir              string
	OutputDir                string
	ClipsDirName             string
	MinClipLengthSeconds     int
	MaxClipLengthSeconds     int
}
