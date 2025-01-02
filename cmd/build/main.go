package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/clipstreamer"
	"gitlab.com/andyinabox/yesterdaysnews/domain/downloader"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

var verbose bool

func init() {
	verbose = false

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	err := godotenv.Load()
	if err != nil {
		log.Warnf("error loading .env: %s", err)
	}

}

func main() {
	var err error
	var dl domain.Downloader
	var cs domain.ClipStreamer

	dl = downloader.New(&downloader.Config{
		GoogleAPIKey: os.Getenv("YN_GOOGLE_API_KEY"),
	})

	downloadDir := "download"
	err = os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	outputDir := "output"
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	cs = clipstreamer.New(dl, &clipstreamer.Config{
		VideoDate:                util.Yesterday(),
		DownloadCountPerPlaylist: 10,
		DownloadDir:              downloadDir,
		OutputDir:                outputDir,
	})

	ctx, _ := context.WithCancelCause(context.Background())

	errHandler := func(err domain.StreamErr) {
		switch err.Type() {
		case domain.StreamErrTODO:
			log.Error("TODO stream error: %s", err)
		default:
			log.Error("Misc stream error: %s", err)
		}
	}

	// TODO: let's just use the playlist IDs and skip this step
	channelNames := []string{"@cnn", "@msnbc", "@foxnews"}
	playlistIDs := make([]string, len(channelNames))

	for i, name := range channelNames {
		playlistId, err := dl.GetChannelPlaylistID(ctx, name)
		if err != nil {
			log.Fatal(err)
		}
		playlistIDs[i] = playlistId
	}

	errs := cs.ErrorStream(ctx, errHandler)
	ids := cs.VideoIDStream(ctx, errs, playlistIDs...)

	for id := range ids {
		fmt.Println(id)
	}

	log.Info("Done")

}
