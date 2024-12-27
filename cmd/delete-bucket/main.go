package main

import (
	"context"
	"flag"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/objectstoreclient"
)

var verbose, forceDelete bool
var bucketName string

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&bucketName, "n", "", "bucket name")
	flag.Parse()

	if bucketName == "" {
		log.Fatal("must provide a bucket name")
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {

	osclient := objectstoreclient.New(&objectstoreclient.Config{
		Endpoint:  os.Getenv("YN_S3_ENDPOINT"),
		AccessKey: os.Getenv("YN_S3_ACCESS_KEY"),
		SecretKey: os.Getenv("YN_S3_SECRET_ACCESS_KEY"),
	})

	err := osclient.DeleteBucket(context.Background(), bucketName, true)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("successfully deleted bucket %q", bucketName)
}
