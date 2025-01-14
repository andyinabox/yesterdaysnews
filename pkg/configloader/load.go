package configloader

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codingconcepts/env"
)

func Load(c any) error {
	return env.Set(c)
}

func LoadJSON(c any, data []byte) error {
	err := json.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("error unmarshaling config data: %w", err)
	}

	return Load(c)
}

func LoadJSONFile(c any, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read config file %q: %w", path, err)
	}
	return LoadJSON(c, data)
}
