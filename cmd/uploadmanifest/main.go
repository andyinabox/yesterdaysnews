package main

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/objectstoreclient"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(log.DebugLevel)
}

func main() {
	client := objectstoreclient.New(&objectstoreclient.Config{
		Endpoint:  os.Getenv("YN_S3_ENDPOINT"),
		AccessKey: os.Getenv("YN_S3_ACCESS_KEY"),
		SecretKey: os.Getenv("YN_S3_SECRET_ACCESS_KEY"),
	})

	key, err := client.UploadFile(context.Background(), "yesterdaysnews-test", "manifest.json", "dist/manifest.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Info(key)
}
