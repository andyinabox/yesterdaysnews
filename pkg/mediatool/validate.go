package mediatool

import (
	"context"
	"fmt"

	"code.andydayton.com/andy/yesterdaysnews/pkg/shellargs"
)

func (t *Tool) Validate(ctx context.Context, path string) error {
	options := shellargs.New().
		AddKeyed("-v", "error").
		Add(path)

	result, err := t.executeFfprobe(ctx, options)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	resultStr := string(result)
	if resultStr != "" {
		return fmt.Errorf("validation error: %s", resultStr)
	}

	return nil
}
