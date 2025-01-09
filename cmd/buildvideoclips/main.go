package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/domain/builder"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

var verbose bool

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}

	err := godotenv.Load()
	if err != nil {
		log.Warnf("error loading .env: %s", err)
	}
}

func main() {
	var b domain.Builder
	var eh domain.ErrorHandler

	ctx := context.Background()
	eh = domain.DefaultErrorHandler(ctx)

	// error recovery
	defer func() {
		// print error report if there are any errors
		if eh.CountAll() > 0 {
			report := eh.Report()
			log.Print(eh.Report())
			_ = os.WriteFile(fmt.Sprintf("errors.%s.json", util.Timestamp(time.Now())), []byte(report), os.ModePerm)
		}
	}()

	// load config
	config := &builder.Config{}
	err := config.Load("builder.config.json")
	if err != nil {
		log.Fatal(err)
	}

	// for testing setting this to 1
	config.PlaylistIDs = []string{"UUupvZG-5ko_eiXAupbDfxWw"}
	config.DownloadCountPerPlaylist = 1

	// making output dirs
	err = os.MkdirAll(filepath.Join(config.OutputDir, config.ClipsDirName), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	b = builder.New(config, eh)

	yesterday := util.Yesterday()
	uploadDir := util.Timestamp(yesterday)

	clips := b.BuildVideoClips(ctx, util.Yesterday(), uploadDir)

	log.Infof("finished processing %d clips", len(clips))
}

// var verbose, outputManifest bool
// var channelName, outputDir string
// var maxResults int

// func init() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	flag.StringVar(&channelName, "n", "@CNN", "username/handle for channel")
// 	flag.StringVar(&outputDir, "o", "output/downloads/cnn", "where to download files")
// 	flag.IntVar(&maxResults, "c", 10, "max number of videos to download")
// 	flag.BoolVar(&verbose, "v", false, "verbose output")
// 	flag.BoolVar(&outputManifest, "m", true, "output manifest file")
// 	flag.Parse()

// 	if verbose {
// 		log.SetLevel(log.DebugLevel)
// 		log.SetReportCaller(true)
// 	}

// }

// func main() {

// 	var dl domain.Downloader

// 	dl = downloader.New(&downloader.Config{
// 		GoogleAPIKey: os.Getenv("GOOGLE_API_KEY"),
// 		BinPathYTDLP: os.Getenv("YT_DLP_PATH"),
// 	})

// 	err := os.MkdirAll(outputDir, os.ModePerm)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	result, err := dl.DownloadVideosForChannel(context.Background(), domain.DownloadRequest{
// 		ChannelUsername: channelName,
// 		Date:            time.Now().AddDate(0, 0, -1),
// 		MaxResults:      maxResults,
// 		OutputDir:       outputDir,
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if outputManifest {
// 		outFile := filepath.Join(outputDir, "manifest.json")

// 		b, err := json.Marshal(result)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		log.Info("outputting manifest file", "file", outFile)
// 		err = os.WriteFile(outFile, b, os.ModePerm)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	log.Infof("finished downloading %d videos", len(result.Files))
// }
