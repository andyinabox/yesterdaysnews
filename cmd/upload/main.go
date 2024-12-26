package main

import (
	"context"
	"flag"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/uploader"
)

var verbose bool

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	flag.BoolVar(&verbose, "v", true, "verbose output")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {

	u := uploader.New(&uploader.Config{
		S3Endpoint:  os.Getenv("YN_S3_ENDPOINT"),
		S3AccessKey: os.Getenv("YN_S3_ACCESS_KEY"),
		S3SecretKey: os.Getenv("YN_S3_SECRET_ACCESS_KEY"),

		BucketNameBase: os.Getenv("YN_S3_BUCKET_NAME"),

		VideoFileName: "yesterdays-news.mp4",
		ModelFileName: "yesterdays-news.model.json",
		ClipsDirName:  "clips",
	})

	containerName, err := u.Upload(context.Background(), "dist")
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("successfully deployed container %q", containerName)

	demotedContainer, err := u.PromoteContainer(context.Background(), containerName)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("successfully promoted container %q", containerName)

	deletedContainers, err := u.PruneContainers(context.Background(), demotedContainer)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("successfully pruned %d containers", len(deletedContainers))

}
