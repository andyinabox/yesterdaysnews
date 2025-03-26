package shell

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

func (s *Shell) GetBinaryPath(ctx context.Context, cmd string) (path string, err error) {
	result, err := s.Execute(ctx, fmt.Sprintf("which %s", cmd))

	if err != nil {
		return
	}

	path = strings.TrimSpace(string(result))

	if path == "" {
		err = errors.New("no path returned from 'which' command")
	}

	return
}

func (s *Shell) MustGetBinaryPath(cmd string) (path string) {
	var err error
	path, err = s.GetBinaryPath(context.Background(), cmd)
	if err != nil {
		panic(fmt.Sprintf("uable to find find binary path", "cmd", cmd, "error", err))
	}
	return
}
