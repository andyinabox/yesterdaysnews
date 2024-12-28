package main

import (
	"context"
	"flag"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain/uploader"
)

var verbose bool

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	flag.BoolVar(&verbose, "v", false, "verbose output")
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

		ContainerName: os.Getenv("YN_S3_BUCKET_NAME"),
		PrimaryDir:    "current",
	})

	log.Info("uploading assets...")

	prefix, err := u.Upload(context.Background(), "dist")
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("successfully deployed objects to %q", prefix)

	demotedPrefix, err := u.PromoteObjects(context.Background(), prefix)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("successfully promoted objects with prefix %q and demoted %q", prefix, demotedPrefix)

	deletedObjects, err := u.PruneObjects(context.Background(), demotedPrefix)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("successfully pruned %d objects", len(deletedObjects))

}
