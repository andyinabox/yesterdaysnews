package youtubedownloader

import (
	"testing"

	"github.com/charmbracelet/log"
)

func TestDownloadVideoRequest(t *testing.T) {

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
