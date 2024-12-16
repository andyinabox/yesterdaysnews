package shell

import (
	"context"
	"strings"

	"github.com/charmbracelet/log"
)

type Shell struct{}

func New() *Shell {
	return &Shell{}
}

// logWriter writes command output to the logger
type logWriter struct {
	ctx context.Context
}

func (l *logWriter) Write(d []byte) (int, error) {
	lines := strings.Split(string(d), "\n")
	for _, l := range lines {
		log.Debug(l)
	}
	return len(d), nil
}

// logWriter writes command output to an output channel
type outputChannelWriter struct {
	output chan<- string
}

func (w *outputChannelWriter) Write(d []byte) (int, error) {
	lines := strings.Split(string(d), "\n")
	for _, l := range lines {
		w.output <- l
	}
	return len(d), nil
}

func (s *Shell) Execute(ctx context.Context, cmd string) ([]byte, error) {
	return s.execute(ctx, cmd, "", nil)
}
func (s *Shell) ExecuteWithInput(ctx context.Context, cmd string, input string) ([]byte, error) {
	return s.execute(ctx, cmd, input, nil)

}
func (s *Shell) ExecuteWithOutput(ctx context.Context, cmd string, output chan<- string) ([]byte, error) {
	return s.execute(ctx, cmd, "", output)
}
func (s *Shell) ExecuteWithInputAndOutput(ctx context.Context, cmd string, input string, output chan<- string) ([]byte, error) {
	return s.execute(ctx, cmd, input, output)
}
