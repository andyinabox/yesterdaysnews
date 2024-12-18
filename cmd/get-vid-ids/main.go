package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdays-news-downloader/domain/downloader"
)

var verbose bool
var userName string
var videoCount int

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&userName, "u", "@CNN", "channel username")
	flag.IntVar(&videoCount, "c", 10, "how many ids to go for")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}
}

func main() {

	dl := downloader.New(&downloader.Config{
		GoogleAPIKey: os.Getenv("GOOGLE_API_KEY"),
		BinPathYTDLP: os.Getenv("YT_DLP_PATH"),
	})

	id, err := dl.GetUploadsPlaylistIdForChannel(context.Background(), userName)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("got user id", "id", id)

	yesterday := time.Now().Add(-24 * time.Hour)

	ids, err := dl.GetPlaylistVideosForDate(context.Background(), id, yesterday, videoCount)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("got ids", "ids", ids)

	for _, i := range ids {
		fmt.Println(i)
	}

}
