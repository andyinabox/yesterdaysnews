package captionschain

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseVTTFileDecodesHTMLEntities(t *testing.T) {
	vtt := "WEBVTT\n\n" +
		"00:00:01.000 --> 00:00:03.000\n" +
		"it&#39;s &gt; all &lt; the &amp; things\n\n"

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.vtt")
	if err := os.WriteFile(path, []byte(vtt), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	c := New(2)
	out, err := c.parseVTTFile(path)
	if err != nil {
		t.Fatalf("parseVTTFile: %v", err)
	}

	want := "it's > all < the & things"
	if !strings.Contains(out, want) {
		t.Errorf("output missing decoded entities\nwant contains: %q\ngot: %q", want, out)
	}

	for _, entity := range []string{"&gt;", "&lt;", "&amp;", "&#39;"} {
		if strings.Contains(out, entity) {
			t.Errorf("output still contains raw entity %q: %q", entity, out)
		}
	}
}
