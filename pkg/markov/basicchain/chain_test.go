package basicchain

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/log"
)

func TestBasicChainJSONMarshaling(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)

	f, err := os.Open("../../../test/corpus.txt")
	if err != nil {
		t.Fatal(err)
	}

	chain1 := *New(2)
	chain1.Build(f)

	data1, err := chain1.Save()
	if err != nil {
		t.Fatal(err)
	}

	if string(data1) == "{}" {
		t.Error("json data is empty")
	}

	chain2 := *New(2)
	err = chain2.Load(data1)
	if err != nil {
		t.Fatal(err)
	}

	data2, err := chain1.Save()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data1, data2) {
		t.Errorf("expected identical json output:\n%s\n%s", string(data1), string(data2))
	}

	if len(chain1.chain) != len(chain2.chain) {
		t.Errorf("expected identical lengths for chains: %d, %d", len(chain1.chain), len(chain2.chain))
	}

	if len(chain1.prefixes) != len(chain2.prefixes) {
		t.Errorf("expected identical lengths for prefixes: %d, %d", len(chain1.prefixes), len(chain2.prefixes))
	}

	if len(chain1.startPrefixes) != len(chain2.startPrefixes) {
		t.Errorf("expected identical lengths for startPrefixes: %d, %d", len(chain1.startPrefixes), len(chain2.startPrefixes))
	}

	if len(chain1.endChain) != len(chain2.endChain) {
		t.Errorf("expected identical lengths for endChains: %d, %d", len(chain1.endChain), len(chain2.endChain))
	}

}

func TestBasicChainStartAndEndWords(t *testing.T) {
	corpus := "This is a very short text."

	chain := New(2)
	chain.Build(strings.NewReader(corpus))

	if chain.startPrefixes[0] != "This is" {
		t.Error("missing startPrefix")
		t.Log(chain.startPrefixes)
	}

	opts, ok := chain.endChain["very short"]
	if !ok {
		t.Fatal("missing endChain prefix")
		t.Log(chain.endChain)
	}

	if opts[0] != "text." {
		t.Error("missing endChain option")
		t.Log(chain.endChain)
	}

}
