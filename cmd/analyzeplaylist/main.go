// analyzeplaylist runs the builder's playlist filter pipeline against a
// single playlist ID and reports per-reason counts. Use it to evaluate
// whether a candidate YouTube channel is viable as a source for the builder.
//
// Example:
//
//	go run ./cmd/analyzeplaylist -playlist UUupvZG-5ko_eiXAupbDfxWw -pages 3
//
// By default it targets yesterday's date (matching production). Override with
// -date YYYY-MM-DD to probe a specific day.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"code.andydayton.com/andy/yesterdaysnews/domain"
	"code.andydayton.com/andy/yesterdaysnews/domain/errorhandler"
	"code.andydayton.com/andy/yesterdaysnews/domain/logger"
	"code.andydayton.com/andy/yesterdaysnews/domain/youtubeservice"
	"code.andydayton.com/andy/yesterdaysnews/pkg/util"
)

func main() {
	var (
		playlistID string
		channel    string
		dateStr    string
		pages      int
		maxSizeMB  int
		verbose    bool
	)

	flag.StringVar(&playlistID, "playlist", "", "playlist ID to analyze")
	flag.StringVar(&channel, "channel", "", "channel handle (e.g. @USATODAY) — resolves to the channel's uploads playlist")
	flag.StringVar(&dateStr, "date", "", "target date as YYYY-MM-DD (defaults to yesterday)")
	flag.IntVar(&pages, "pages", 3, "maximum number of 50-item pages to fetch")
	flag.IntVar(&maxSizeMB, "max-size", 2048, "max video filesize in MiB")
	flag.BoolVar(&verbose, "v", false, "verbose logging (shows per-item skip reasons)")
	flag.Parse()

	logger.SetDefault(&logger.Config{Verbose: verbose})

	if (playlistID == "") == (channel == "") {
		slog.Error("exactly one of -playlist or -channel is required")
		flag.Usage()
		os.Exit(1)
	}

	if err := godotenv.Load(); err != nil {
		slog.Warn("error loading .env", "error", err)
	}

	targetDate := util.Yesterday()
	if dateStr != "" {
		d, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
		if err != nil {
			slog.Error("invalid -date, expected YYYY-MM-DD", "error", err)
			os.Exit(1)
		}
		targetDate = d
	}

	maxSize := uint(maxSizeMB) * 1024 * 1024

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eh := errorhandler.New(ctx, &errorhandler.Config{})
	errs := eh.Channel()

	yt := youtubeservice.New(&youtubeservice.Config{
		GoogleAPIKey:        os.Getenv("YN_GOOGLE_API_KEY"),
		BinPathYTDLP:        os.Getenv("YN_YT_DLP_PATH"),
		MaxPlaylistRequests: pages,
		ThrottleDownloadsBy: 0,
	})

	if channel != "" {
		resolved, err := yt.GetChannelPlaylistID(ctx, channel)
		if err != nil {
			slog.Error("failed to resolve channel handle", "channel", channel, "error", err)
			os.Exit(1)
		}
		if resolved == "" {
			slog.Error("channel resolved to empty uploads playlist (handle may be wrong)", "channel", channel)
			os.Exit(1)
		}
		fmt.Printf("resolved channel %s → playlist %s\n", channel, resolved)
		playlistID = resolved
	}

	var mu sync.Mutex
	counts := make(map[domain.FilterReason]int)
	onFilter := func(reason domain.FilterReason) {
		mu.Lock()
		counts[reason]++
		mu.Unlock()
	}

	fmt.Printf("analyzing playlist %q for target date %s\n", playlistID, targetDate.Format("2006-01-02"))
	fmt.Printf("max pages: %d  (up to %d items)\n", pages, pages*50)
	fmt.Printf("format filter: %s\n", youtubeservice.VideoFormatString)
	fmt.Println()

	start := time.Now()

	// Pass math.MaxInt so the stream doesn't early-stop after N accepted IDs —
	// we want the full breakdown across every fetched item.
	stream := yt.GetPlaylistVideoIDStream(ctx, errs, playlistID, targetDate, maxSize, math.MaxInt, onFilter)

	var accepted []string
	for id := range stream {
		accepted = append(accepted, id)
	}

	elapsed := time.Since(start)

	printReport(counts, accepted, elapsed)
}

func printReport(counts map[domain.FilterReason]int, accepted []string, elapsed time.Duration) {
	order := []domain.FilterReason{
		domain.FilterReasonFetched,
		domain.FilterReasonWrongDate,
		domain.FilterReasonFormatUnavailable,
		domain.FilterReasonTooLarge,
		domain.FilterReasonNoCaptions,
		domain.FilterReasonInfoError,
		domain.FilterReasonAccepted,
	}

	fetched := counts[domain.FilterReasonFetched]

	fmt.Println("Filter breakdown:")
	fmt.Printf("  %-20s %6s  %6s\n", "reason", "count", "% of fetched")
	for _, r := range order {
		c := counts[r]
		var pct string
		if fetched > 0 {
			pct = fmt.Sprintf("%5.1f%%", float64(c)/float64(fetched)*100)
		} else {
			pct = "   —"
		}
		fmt.Printf("  %-20s %6d  %s\n", r, c, pct)
	}
	fmt.Println()

	// Show items that made it through the date filter but still got rejected,
	// so you can see how much of the "viable pool" each later stage kills.
	postDate := fetched - counts[domain.FilterReasonWrongDate]
	fmt.Printf("Items on target date: %d\n", postDate)
	if postDate > 0 {
		rejected := counts[domain.FilterReasonFormatUnavailable] +
			counts[domain.FilterReasonTooLarge] +
			counts[domain.FilterReasonNoCaptions] +
			counts[domain.FilterReasonInfoError]
		fmt.Printf("  rejected by later filters: %d\n", rejected)
		fmt.Printf("  accepted:                  %d\n", counts[domain.FilterReasonAccepted])
	}
	fmt.Println()

	if len(accepted) > 0 {
		sort.Strings(accepted)
		fmt.Printf("Accepted video IDs (%d):\n", len(accepted))
		for _, id := range accepted {
			fmt.Printf("  https://youtube.com/watch?v=%s\n", id)
		}
	} else {
		fmt.Println("No videos accepted.")
	}

	fmt.Printf("\nElapsed: %s\n", elapsed.Round(time.Millisecond))
}
