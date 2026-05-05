package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/domain/logger"
	"gitlab.com/andyinabox/yesterdaysnews/domain/youtubeservice"
)

var channelName string
var verbose bool

func init() {
	flag.StringVar(&channelName, "n", "", "channel name")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	logger.SetDefault(&logger.Config{
		Verbose: verbose,
	})

	if channelName == "" {
		slog.Error("channel name is required")
		os.Exit(1)
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
		panic(err)
	}

	fmt.Println(id)

}
