package builder

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	envTagName = "env"
)

type Config struct {
	// env variables
	GoogleAPIKey string `env:"GOOGLE_API_KEY"`
	S3Endpoint   string `env:"YN_S3_ENDPOINT"`
	S3AccessKey  string `env:"YN_S3_ACCESS_KEY"`
	S3SecretKey  string `env:"YN_S3_SECRET_ACCESS_KEY"`

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
}

func (c *Config) Load(configFile string) error {

	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("error unmarshaling config file: %w", err)
	}

	// doing this until I get struct tag working
	c.GoogleAPIKey = os.Getenv("GOOGLE_API_KEY")
	c.S3Endpoint = os.Getenv("YN_S3_ENDPOINT")
	c.S3AccessKey = os.Getenv("YN_S3_ACCESS_KEY")
	c.S3SecretKey = os.Getenv("YN_S3_SECRET_ACCESS_KEY")

	return nil
}
