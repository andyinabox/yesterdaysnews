package shell

import (
	"errors"
	"fmt"
	"os/exec"
)

func wrapErr(err error) error {

	if err == nil {
		return nil
	}

	var eerr *exec.Error
	if errors.As(err, &eerr) {
		return fmt.Errorf("shell command cannot be executed: %w", err)
	}

	if exerr, ok := err.(*exec.ExitError); ok {
		switch exerr.ExitCode() {
		case 11:
			return fmt.Errorf("shell command maybe wrong permissions: %w", err)
		case 12:
			return fmt.Errorf("shell command maybe non-existent command: %w", err)
		}
	}

	return fmt.Errorf("shell execution error: %w", err)
}
