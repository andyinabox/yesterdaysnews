package builder

import (
	"os"
	"testing"
)

func TestConfigLoader(t *testing.T) {

	expectedAPIKey := "test"
	expectedMinClipLength := 5

	os.Setenv("GOOGLE_API_KEY", expectedAPIKey)

	conf := &Config{}

	err := conf.Load("../../builder.config.json")
	if err != nil {
		t.Fatal(err)
	}

	if conf.GoogleAPIKey != expectedAPIKey {
		t.Errorf("expected env var to be %q, got %q", expectedAPIKey, conf.GoogleAPIKey)
	}

	if conf.MinClipLengthSeconds != expectedMinClipLength {
		t.Errorf("expected json var to be %d, got %d", expectedMinClipLength, conf.MinClipLengthSeconds)
	}

}
