package basicchain

import (
	"os"
	"testing"

	"github.com/charmbracelet/log"
)

func init() {
	log.SetLevel(log.DebugLevel)
}
func TestLoadV0(t *testing.T) {
	file := "../../../test/model-v0.json"

	chain := New(0)

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	err = chain.Load(data)
	if err != nil {
		t.Fatal(err)
	}

	if len(chain.chain) == 0 {
		t.Error("chain length should not be 0")
	}

	t.Log(chain.Start())
	t.Log(chain.Start())
	t.Log(chain.Start())
}

func TestLoadV1(t *testing.T) {
	file := "../../../test/model-v1.json"

	chain := New(0)

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	err = chain.Load(data)
	if err != nil {
		t.Fatal(err)
	}

	if len(chain.chain) == 0 {
		t.Error("chain length should not be 0")
	}

	t.Log(chain.Start())
	t.Log(chain.Start())
	t.Log(chain.Start())
}
