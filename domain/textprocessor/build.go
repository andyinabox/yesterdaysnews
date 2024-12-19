package textprocessor

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/markov"
)

type CorpusType string

const (
	CorpusTypeText CorpusType = "txt"
	CorpusTypeVTT  CorpusType = "vtt"
)

type Corpus struct {
	Type     CorpusType
	FileGlob string
	Weight   float32
}

func (p *Processor) Build(sources []Corpus) error {

	corpi := make([]markov.MultiChainSub, len(sources))

	// iterate through sources
	for sourceIndex, source := range sources {
		files, err := filepath.Glob(source.FileGlob)
		if err != nil {
			return fmt.Errorf("error resolving file glob: %s: %w", source.FileGlob, err)
		}

		corpus := ""

		// iterate through files
		for _, f := range files {
			var err error
			var parsedContent string

			// parse based on type
			switch source.Type {
			case CorpusTypeText:
				parsedContent, err = p.parseTextFile(f)
			case CorpusTypeVTT:
				parsedContent, err = p.parseVTTFile(f)
			default:
				return fmt.Errorf("invalid CorpusType: %s", source.Type)
			}

			if err != nil {
				return err
			}

			// add to corpus
			corpus = fmt.Sprintf("%s\n%s", corpus, parsedContent)
		}

		chain := markov.NewBasicChain(p.cfg.PrefixLength)
		chain.Build(bytes.NewReader([]byte(corpus)))

		corpi[sourceIndex] = markov.MultiChainSub{
			Chain:  chain,
			Weight: source.Weight,
		}
	}

	return p.chain.Build(corpi)
}

func (p *Processor) parseTextFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file %q: %w", filePath, err)
	}

	return string(content), nil
}

func (p *Processor) parseVTTFile(filePath string) (string, error) {
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
