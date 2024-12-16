package shell

import (
	"context"
	"strings"
	"testing"
)

func TestExecuteWithInput(t *testing.T) {
	sh := New()

	result, err := sh.ExecuteWithInput(context.Background(), "cat", "success")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(result))

	// this will obviously fail on a non-mac
	if !strings.Contains(string(result), "success") {
		t.Fatal("Stdin did not get written correctly")
	}

}
