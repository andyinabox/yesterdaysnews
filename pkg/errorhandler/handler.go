package errorhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

var defaultErrorFunc, defaultFatalFunc func(string, error)

func init() {
	defaultErrorFunc = func(typ string, err error) {
		log.Printf("%s error: %s\n", typ, err)
	}
	defaultFatalFunc = func(typ string, err error) {
		panic(fmt.Errorf("%s error: %w", typ, err))
	}
}

type ErrorHandler interface {
	Add(typ string, err error)
	Channel() chan<- Error
	Count(string) int
	CountAll() int
	Report() string
	DeferredReport()
}

type Config struct {
	ErrorFunc  func(string, error)
	FatalFunc  func(string, error)
	Thresholds map[string]int
}

type errorHandler struct {
	errorFunc  func(string, error)
	fatalFunc  func(string, error)
	thresholds map[string]int
	errs       map[string][]error
	stream     chan Error
}

func New(ctx context.Context, cfg *Config) ErrorHandler {
	var mu sync.Mutex
	// ctx, cancel = context.WithCancelCause(ctx)
	stream := make(chan Error)

	h := &errorHandler{
		errorFunc:  defaultErrorFunc,
		fatalFunc:  defaultFatalFunc,
		thresholds: make(map[string]int),
		errs:       make(map[string][]error),
		stream:     stream,
	}

	if cfg.ErrorFunc != nil {
		h.errorFunc = cfg.ErrorFunc
	}

	if cfg.FatalFunc != nil {
		h.fatalFunc = cfg.FatalFunc
	}

	if cfg.Thresholds != nil {
		h.thresholds = cfg.Thresholds
	}

	handleErr := func(err Error) {

		typ := err.Type()

		// add to errors
		mu.Lock()
		h.errs[typ] = append(h.errs[typ], err)
		mu.Unlock()

		// check threshold
		limit, found := h.thresholds[err.Type()]
		if found && h.Count(typ) > limit {
			h.fatalFunc(typ, fmt.Errorf("recieved %d %q errors, limit is %d: %w", h.Count(typ), typ, limit, err))
			return
		}

		// handle error normally
		h.errorFunc(typ, err)
	}

	// main gorhoutine
	go func() {
		defer close(stream)
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-stream:
				handleErr(err)
			}
		}
	}()

	return h
}

func (h *errorHandler) Add(typ string, err error) {
	h.stream <- Err(typ, err)
}

func (h *errorHandler) Count(typ string) int {
	errs, found := h.errs[typ]
	if !found {
		return 0
	}
	return len(errs)
}

func (h *errorHandler) CountAll() int {
	count := 0

	for _, errs := range h.errs {
		count += len(errs)
	}

	return count
}

func (h *errorHandler) Channel() chan<- Error {
	return h.stream
}

func (h *errorHandler) Report() string {
	data, err := json.MarshalIndent(h.errs, "", "  ")
	if err != nil {
		log.Fatal(fmt.Errorf("unable to generate error report: %w", err))
	}
	return string(data)
}

// DeferredReport will log the error report and output an errors file if there are errors (use like `defer h.DeferredReport()`)
func (h *errorHandler) DeferredReport() {
	if h.CountAll() > 0 {
		report := h.Report()
		log.Print(h.Report())
		_ = os.WriteFile(fmt.Sprintf("errors.%s.json", util.Timestamp(time.Now())), []byte(report), os.ModePerm)
	}
}
