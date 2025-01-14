package configloader

import (
	"fmt"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	config := struct {
		EnvValue string `env:"TEST_VALUE"`
	}{}

	expected := "ok"
	os.Setenv("TEST_VALUE", expected)

	err := Load(&config)
	if err != nil {
		t.Fatal(err)
	}

	if config.EnvValue != expected {
		t.Errorf("expected %q, got %q", expected, config.EnvValue)
	}
}

func TestLoadJSON(t *testing.T) {
	expectedConfig := "success"
	jsonData := fmt.Sprintf(`{ "ConfigValue" : %q }`, expectedConfig)

	expectedEnv := "ok"
	os.Setenv("TEST_VALUE", expectedEnv)

	config := struct {
		EnvValue    string `env:"TEST_VALUE"`
		ConfigValue string
	}{}

	err := LoadJSON(&config, []byte(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	if config.EnvValue != expectedEnv {
		t.Errorf("expected %q, got %q", expectedEnv, config.EnvValue)
	}

	if config.ConfigValue != expectedConfig {
		t.Errorf("expected %q, got %q", expectedConfig, config.ConfigValue)
	}
}

func TestLoadJSONFile(t *testing.T) {
	expectedConfig := "success"
	jsonData := fmt.Sprintf(`{ "ConfigValue" : %q }`, expectedConfig)
	jsonFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	jsonFile.WriteString(jsonData)
	jsonFile.Close()

	t.Cleanup(func() {
		os.Remove(jsonFile.Name())
	})

	expectedEnv := "ok"
	os.Setenv("TEST_VALUE", expectedEnv)

	config := struct {
		EnvValue    string `env:"TEST_VALUE"`
		ConfigValue string
	}{}

	err = LoadJSONFile(&config, jsonFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if config.EnvValue != expectedEnv {
		t.Errorf("expected %q, got %q", expectedEnv, config.EnvValue)
	}

	if config.ConfigValue != expectedConfig {
		t.Errorf("expected %q, got %q", expectedConfig, config.ConfigValue)
	}
}
