package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/objectstoreclient"
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

	c := objectstoreclient.New(&objectstoreclient.Config{
		Endpoint:  os.Getenv("YN_S3_ENDPOINT"),
		AccessKey: os.Getenv("YN_S3_ACCESS_KEY"),
		SecretKey: os.Getenv("YN_S3_SECRET_ACCESS_KEY"),
	})

	prefixes, err := c.ListPrefixes(context.Background(), os.Getenv("YN_S3_BUCKET_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	for _, prefix := range prefixes {
		fmt.Println(prefix)
	}
}
