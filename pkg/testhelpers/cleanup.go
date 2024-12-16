package testhelpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileCleanup(t *testing.T, fileGlob string) func() {
	return func() {
		files, err := filepath.Glob(filepath.Join("test/output", fileGlob))
		if err != nil {
			t.Error(err)
		}
		for _, f := range files {
			err = os.Remove(f)
			if err != nil {
				t.Error(err)
			}
		}
	}
}
