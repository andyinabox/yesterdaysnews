package youtubedownloader

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/charmbracelet/log"
)

func TestDownloadVideoRequestToString(t *testing.T) {

	// test all
	{
		req := &Request{
			// Format
			Format: "bv",

			// Subtitles
			WriteAutoSubs: true,
			SubFormat:     "vtt",

			// Verbosity & simulateion
			Quiet:          true,
			NoSimulate:     true,
			DumpJSON:       true,
			DumpSingleJSON: true,

			// File outpuyt
			Output: "file.mp4",
		}

		expected := " --format 'bv' --write-auto-subs --sub-format 'vtt' --quiet --no-simulate --dump-json --dump-single-json --output 'file.mp4'"

		result := req.String()

		if result != expected {
			log.Errorf("expected:\n%q\ngot:\n%q\n", expected, result)
		}

	}

	// test none
	{
		req := &Request{}

		expected := ""

		result := req.String()

		if result != expected {
			log.Errorf("expected:\n%q\ngot:\n%q\n", expected, result)
		}

	}
}

func TestDownloadVideoRequestUnmarshall(t *testing.T) {
	data, err := os.ReadFile("../../test/video-info.json")
	if err != nil {
		t.Fatal(err)
	}

	ytv := YouTubeVideo{}
	err = json.Unmarshal(data, &ytv)
	if err != nil {
		t.Fatal(err)
	}

	expectedId := "8N5yiQ1SABE"

	if ytv.ID != expectedId {
		t.Errorf("expected id %q, got %q", expectedId, ytv.ID)
	}
}
