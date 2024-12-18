package textprocessor

import (
	"bufio"
	"context"
	"os"
	"strings"
	"unicode"

	"github.com/charmbracelet/log"
)

func (p *Processor) VTTToText(ctx context.Context, filePath string) (string, error) {
	log.Debug("parse vtt file", "path", filePath)
	file, err := os.Open(filePath)
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
