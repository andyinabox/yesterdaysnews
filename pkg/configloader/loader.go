package configloader

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

const envTagName = "env"

func Load(pointer any, configFile string) error {
	err := godotenv.Load()
	if err != nil {
		log.Warnf("error loading .env: %s", err)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	err = json.Unmarshal(data, pointer)
	if err != nil {
		return fmt.Errorf("error unmarshaling config file: %w", err)
	}

	t := reflect.TypeOf(pointer)
	v := reflect.ValueOf(pointer)

	for i := 0; i < t.NumField(); i++ {

		// get type and value for field
		typeField := t.Field(i)
		valField := v.Field(i)
		// Get the field tag value
		tag := typeField.Tag.Get(envTagName)
		valField.SetString(os.Getenv(strings.TrimSpace(tag)))
	}

	return nil
}
