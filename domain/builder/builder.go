package builder

import "time"

type Config struct {
	GoogleAPIKey             string
	S3Endpoint               string
	S3AccessKey              string
	S3SecretKey              string
	ObjectStoreContainerName string
	ObjectStorePrimaryDir    string
	DownloadDir              string
	OutputDir                string
	ClipsDirName             string
	MinClipLength            time.Duration
	MaxClipLength            time.Duration
}
