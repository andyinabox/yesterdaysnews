package logger

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

type Config struct {
	Verbose    bool
	JSONOutput bool
	WithAttr   []any
}

func New(c *Config) (logger *slog.Logger) {

	// for json output we'll just use the slog library
	if c.JSONOutput {
		logger = newSlog(c)
		// for text output we'll use charmbracelet
	} else {
		logger = newCharm(c)
	}

	if c.WithAttr != nil {
		logger = logger.With(c.WithAttr...)
	}

	return
}

func newSlog(c *Config) (logger *slog.Logger) {
	var options *slog.HandlerOptions
	if c.Verbose {
		options = &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
	}

	logger = slog.New(slog.NewJSONHandler(os.Stdout, options))

	return
}

func newCharm(c *Config) (logger *slog.Logger) {
	handler := log.New(os.Stderr)

	if c.Verbose {
		handler.SetLevel(log.DebugLevel)
		handler.SetReportCaller(true)
	}

	logger = slog.New(handler)
	return
}

func SetDefault(c *Config) {
	logger := New(c)
	slog.SetDefault(logger)
}
