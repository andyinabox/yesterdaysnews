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
	"gitlab.com/andyinabox/yesterdaysnews/pkg/errorhandler"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

var verbose bool

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	log.SetReportCaller(true)
	log.SetReportTimestamp(false)

	if verbose {
		log.SetLevel(log.DebugLevel)
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

	// set up error handling
	fatalFunc := func(typ string, err error) {
		log.Fatalf("%s: %s", typ, err)
	}
	errorFunc := func(typ string, err error) {
		if typ == domain.ErrTypeFatal {
			fatalFunc(typ, err)
		}
		log.Errorf("%s: %s", typ, err)
	}
	eh = errorhandler.New(ctx, &errorhandler.Config{
		ErrorFunc: errorFunc,
		FatalFunc: fatalFunc,
	})

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
	config.DownloadCountPerPlaylist = 1

	// making output dirs
	err = os.MkdirAll(filepath.Join(config.OutputDir, config.ClipsDirName), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	b = builder.New(config, eh)

	err = b.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
