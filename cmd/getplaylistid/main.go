package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain/youtubeservice"
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
	dl := youtubeservice.New(&youtubeservice.Config{
		GoogleAPIKey:        os.Getenv("YN_GOOGLE_API_KEY"),
		MaxPlaylistRequests: 10,
		ThrottleDownloadsBy: time.Duration(0),
	})

	id, err := dl.GetChannelPlaylistID(context.Background(), channelName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(id)

}
