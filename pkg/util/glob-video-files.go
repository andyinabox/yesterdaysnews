package util

import "path/filepath"

func GlobVideoFiles(dir, base string) []string {
	webm, _ := filepath.Glob(filepath.Join(dir, base+".webm"))
	mp4, _ := filepath.Glob(filepath.Join(dir, base+".mp4"))
	return append(webm, mp4...)
}
