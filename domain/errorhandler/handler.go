package errorhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

var defaultErrorFunc, defaultFatalFunc func(string, error)

func init() {
	defaultErrorFunc = func(typ string, err error) {
		log.Printf("%s error: %s\n", typ, err)
	}
}

type Config struct {
	ErrorFunc  func(string, error)
	Thresholds map[string]int
}

type errorHandler struct {
	errorFunc  func(string, error)
	fatalFunc  func(string, error)
	thresholds map[string]int
	errs       map[string][]domain.Error
	stream     chan domain.Error
	err        error
}

func New(ctx context.Context, cfg *Config) domain.ErrorHandler {
	var mu sync.Mutex
	// ctx, cancel = context.WithCancelCause(ctx)
	stream := make(chan domain.Error)

	h := &errorHandler{
		errorFunc:  defaultErrorFunc,
		thresholds: make(map[string]int),
		errs:       make(map[string][]domain.Error),
		stream:     stream,
	}

	if cfg.ErrorFunc != nil {
		h.errorFunc = cfg.ErrorFunc
	}

	if cfg.Thresholds != nil {
		h.thresholds = cfg.Thresholds
	}

	handleErr := func(err domain.Error) {

		typ := err.Type()

		// add to errors
		mu.Lock()
		h.errs[typ] = append(h.errs[typ], err)
		mu.Unlock()

		// check threshold
		limit, found := h.thresholds[err.Type()]
		if found && h.Count(typ) > limit {
			h.err = fmt.Errorf("recieved %d %q errors, limit is %d: %w", h.Count(typ), typ, limit, h.err)
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

func (h *errorHandler) Channel() chan<- domain.Error {
	return h.stream
}

func (h *errorHandler) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.errs)
}

func (h *errorHandler) String() (str string) {

	for typ, errs := range h.errs {
		str += fmt.Sprintf("%s:\n", typ)
		for i, err := range errs {
			str += fmt.Sprintf("  %d: %q\n", i+1, err.Error())
		}
	}

	return
}

// Report will log the error report and output an errors file if there are errors (use like `defer h.DeferredReport()`)
func (h *errorHandler) Report() {
	if h.CountAll() > 0 {

		log.Println(h.String())

		// don't output file if no filename is provided
		data, err := h.MarshalJSON()
		if err != nil {
			log.Printf("error marshaling error data: %s\n", err)
			return
		}

		_ = os.WriteFile(fmt.Sprintf("errors.%s.json", util.Timestamp(time.Now())), data, os.ModePerm)

	} else {
		log.Println("no errors to report")
	}
}

// func (h *errorHandler) Reset() {
// 	h.err = nil
// }

// func (h *errorHandler) Err() error {
// 	return h.err
// }
