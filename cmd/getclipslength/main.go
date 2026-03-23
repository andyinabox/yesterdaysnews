package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/domain/logger"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/mediatool"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

var channelName string
var verbose bool

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	logger.SetDefault(&logger.Config{
		Verbose: verbose,
	})

	err := godotenv.Load()
	if err != nil {
		slog.Warn("error getting .env", "error", err)
	}
}

func main() {

	ctx := context.Background()

	clips, err := util.GlobVideoFiles("dist/clips", "*")
	if err != nil {
		slog.Error("error globbing clip files", "error", err)
		return
	}

	mt := mediatool.New("", "")

	var total mediatool.Duration
	for _, fn := range clips {
		length, err := mt.GetVideoLength(ctx, fn)
		if err != nil {
			slog.Error("error getting video length", "file", fn, "error", err)
		}

		fmt.Printf("%s: %s\n", fn, length)
		total = total + length
	}

	fmt.Printf("TOTAL: %s\n", total)

}
