package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/charmbracelet/log"
)

var verbose bool
var glob string

func init() {
	flag.StringVar(&glob, "f", "", "glob for files to process")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	if glob == "" {
		log.Fatal("no glob")
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {

	files, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("parse glob", "glob", glob, "files", files)

	corpus := ""

	for _, f := range files {
		text, err := vttToText(f)
		if err != nil {
			log.Error(err)
			continue
		}

		corpus = corpus + "\n" + text
	}

	fmt.Print(corpus)

}

func vttToText(path string) (string, error) {
	log.Debug("parse vtt file", "path", path)
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var startInput bool
	lines := []string{}
	scanner := bufio.NewScanner(file)

scanloop:
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// skip empty lines
		if line == "" {
			continue
		}

		if unicode.IsDigit([]rune(line)[0]) {
			// don't start input until we've recieved our first timecode line
			startInput = true
			continue scanloop
		}

		// skip liens with timecode markup included with text
		if strings.Contains(line, "<c>") {
			continue scanloop
		}

		if startInput {
			// avoid duplicates
			for _, s := range lines {
				if line == s {
					continue scanloop
				}
			}
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	output := strings.Join(lines, " ")

	log.Debug(output)

	return output, nil
}
