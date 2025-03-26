package logger

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

type Config struct {
	Verbose    bool
	JSONOutput bool
}

func New(c *Config) (logger *slog.Logger) {

	// for json output we'll just use the slog library
	if c.JSONOutput {
		return newSlog(c)
	}

	// for text output we'll use charmbracelet
	return newCharm(c)
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
