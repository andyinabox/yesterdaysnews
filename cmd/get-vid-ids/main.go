package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/ytapi"
)

var verbose bool
var userName string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&userName, "u", "@CNN", "channel username")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	apiClient := ytapi.New(os.Getenv("GOOGLE_API_KEY"))

	id, err := apiClient.GetUploadsPlaylistIdForUser(context.Background(), "@CNN")
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("got user id", "id", id)

	ids, err := apiClient.GetPlaylistVideosForDate(context.Background(), id, time.Now().Add(-24*time.Hour), 10)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("got ids", "ids", ids)

	for _, i := range ids {
		fmt.Println(i)
	}

}
