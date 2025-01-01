package util

import (
	"errors"
	"os"
)

func DoesFileExist(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}
