package errorhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
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
	Add(Error)
	Count(string) int
	CountAll() int
	Context() context.Context
	Channel() chan<- Error
	Report() string
}

type Config struct {
	ErrorFunc  func(string, error)
	FatalFunc  func(string, error)
	Thresholds map[string]int
}

type errorHandler struct {
	ctx        context.Context
	errorFunc  func(string, error)
	fatalFunc  func(string, error)
	thresholds map[string]int
	errs       map[string][]Error
	stream     chan Error
}

func New(ctx context.Context, cfg *Config) ErrorHandler {
	var mu sync.Mutex
	// ctx, cancel = context.WithCancelCause(ctx)
	stream := make(chan Error)

	h := &errorHandler{
		ctx:        ctx,
		errorFunc:  defaultErrorFunc,
		fatalFunc:  defaultFatalFunc,
		thresholds: make(map[string]int),
		errs:       make(map[string][]Error),
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

func (h *errorHandler) Add(err Error) {
	h.stream <- err
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
func (h *errorHandler) Context() context.Context {
	return h.ctx
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
