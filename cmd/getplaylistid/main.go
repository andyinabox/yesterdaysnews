package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain/downloader"
)

var channelName string
var verbose bool

func init() {
	flag.StringVar(&channelName, "n", "", "channel name")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if channelName == "" {
		log.Fatal("channel name is required")
	}

	err := godotenv.Load()
	if err != nil {
		log.Warnf("error getting .env: %s", err)
	}
}

func main() {
	dl := downloader.New(&downloader.Config{
		GoogleAPIKey: os.Getenv("GOOGLE_API_KEY"),
	})

	id, err := dl.GetChannelPlaylistID(context.Background(), channelName)
	if err != nil {
		log.Fatal(id)
	}

	fmt.Println(id)

}
