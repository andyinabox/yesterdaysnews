package main

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/uploader"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(log.DebugLevel)
}

func main() {

	u := uploader.New(&uploader.Config{
		S3Endpoint:     os.Getenv("YN_S3_ENDPOINT"),
		S3AccessKey:    os.Getenv("YN_S3_ACCESS_KEY"),
		S3SecretKey:    os.Getenv("YN_S3_SECRET_ACCESS_KEY"),
		BucketNameBase: os.Getenv("YN_S3_BUCKET_NAME"),
	})

	deleted, err := u.PruneContainers(context.Background(), "")
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("deleted %d containers", len(deleted))

}
