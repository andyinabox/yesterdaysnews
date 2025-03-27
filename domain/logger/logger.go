package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
	"github.com/grafana/loki-client-go/loki"
	slogloki "github.com/samber/slog-loki/v3"
)

type LoggerType string

const (
	LoggerTypeText LoggerType = "text"
	LoggerTypeJSON LoggerType = "json"
	LoggerTypeLoki LoggerType = "loki"
)

type Config struct {
	Type         LoggerType
	Verbose      bool
	WithAttr     []any
	LokiEndpoint string
	LokiTenentID string
}

// New creates a new `slog.Logger`
func New(c *Config) (logger *slog.Logger) {

	// use text logger by default
	if c.Type == "" {
		c.Type = LoggerTypeText
	}

	switch c.Type {
	case LoggerTypeText:
		logger = newText(c)
	case LoggerTypeJSON:
		logger = newJson(c)
	case LoggerTypeLoki:
		var err error
		logger, err = newLoki(c)
		if err != nil {
			panic(fmt.Errorf("error starting loki logger: %w", err))
		}
	default:
		panic(fmt.Errorf("invalid logger type: %s", c.Type))
	}

	if c.WithAttr != nil {
		logger = logger.With(c.WithAttr...)
	}

	return
}

// SetDefault creates a new `slog.Logger` and sets it as the default
func SetDefault(c *Config) {
	logger := New(c)
	slog.SetDefault(logger)
}

func newLoki(c *Config) (logger *slog.Logger, err error) {
	// setup loki client
	var config loki.Config
	config, err = loki.NewDefaultConfig(c.LokiEndpoint)
	if err != nil {
		err = fmt.Errorf("error creating loki config: %w", err)
		return
	}
	config.TenantID = c.LokiTenentID

	var client *loki.Client
	client, err = loki.New(config)
	if err != nil {
		err = fmt.Errorf("error creating loki config: %w", err)
		return
	}

	options := slogloki.Option{
		Client: client,
	}

	if c.Verbose {
		options.Level = slog.LevelDebug
	}

	logger = slog.New(options.NewLokiHandler())

	return
}

func newJson(c *Config) (logger *slog.Logger) {
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

func newText(c *Config) (logger *slog.Logger) {
	handler := log.New(os.Stderr)

	if c.Verbose {
		handler.SetLevel(log.DebugLevel)
		handler.SetReportCaller(true)
	}

	logger = slog.New(handler)
	return
}
