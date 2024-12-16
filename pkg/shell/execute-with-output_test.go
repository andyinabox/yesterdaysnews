package shell

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestExecuteWithOutput(t *testing.T) {
	sh := New()

	outputStream := make(chan string)
	outputs := []string{}
	ctx, cancel := context.WithCancel(context.Background())

	// capture output lines
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case o := <-outputStream:
				t.Log(o)
				outputs = append(outputs, o)
			}
		}
	}()

	// cancel after 5 seconds
	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()

	result, err := sh.ExecuteWithOutput(ctx, "ping go.dev", outputStream)
	// note that currently we will still get an error if the process is
	// cancelled using the context, but that's expected
	if err != nil && !strings.Contains(err.Error(), "killed") {
		t.Fatal(err)
	}

	t.Log(string(result))

	if len(outputs) == 0 {
		t.Error("no output recieved")
	}
}
