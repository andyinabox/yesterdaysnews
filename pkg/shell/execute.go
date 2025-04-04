package shell

import (
	"bytes"
	"context"
	"io"
	"os/exec"
)

func (s *Shell) execute(ctx context.Context, command string, input string, output chan<- string) (result []byte, err error) {

	// slog.Debug("shell execute", "cmd", command)

	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	var buf bytes.Buffer
	logger := &logWriter{ctx}
	writers := []io.Writer{logger, &buf}

	if output != nil {
		writers = append(writers, &outputChannelWriter{output})
	}

	cmd.Stdout = io.MultiWriter(writers...)
	cmd.Stderr = io.MultiWriter(writers...)

	if input != "" {
		var stdin io.WriteCloser
		stdin, err = cmd.StdinPipe()
		if err != nil {
			return nil, wrapErr(err)
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, input)
		}()

		err = cmd.Start()
		if err != nil {
			return nil, wrapErr(err)
		}

		err = cmd.Wait()

		result = buf.Bytes()
	} else {
		err = cmd.Run()
		result = buf.Bytes()
	}

	// sometimes the result is helpful in handling the error
	return result, wrapErr(err)
}
