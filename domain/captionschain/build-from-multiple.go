package captionschain

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (c *Chain) BuildFromMultiple(sources []domain.Corpus) error {

	combined := ""

	// iterate through sources
	for _, source := range sources {
		log.Debugf("processing text from source %q", source.FileGlob)
		files, err := filepath.Glob(source.FileGlob)
		if err != nil {
			return fmt.Errorf("error resolving file glob: %s: %w", source.FileGlob, err)
		}

		corpus := ""

		// iterate through files
		for _, f := range files {
			log.Debugf("processing file %q", f)
			var err error
			var parsedContent string

			// parse based on type
			switch source.Type {
			case domain.CorpusTypeText:
				parsedContent, err = c.parseTextFile(f)
			case domain.CorpusTypeVTT:
				parsedContent, err = c.parseVTTFile(f)
			default:
				return fmt.Errorf("invalid CorpusType: %s", source.Type)
			}

			if err != nil {
				return err
			}

			// add to corpus
			corpus = fmt.Sprintf("%s\n%s", corpus, parsedContent)
		}

		for i := 0; i < source.Weight; i++ {
			combined = combined + "\n" + corpus
		}

	}

	log.Debug(combined)

	c.BasicChain.Build(strings.NewReader(combined))

	return nil
}

func (c *Chain) parseTextFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file %q: %w", filePath, err)
	}

	return string(content), nil
}

func (c *Chain) parseVTTFile(filePath string) (string, error) {
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
