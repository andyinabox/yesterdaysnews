package builder

import (
	"math"
	"testing"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func newTestBuilder(playlists []domain.PlaylistSource) *Builder {
	return &Builder{
		cfg:              &Config{Playlists: playlists},
		videoSourceMap:   make(map[string]string),
		sourceDurations:  make(map[string]time.Duration),
		sourceClipCounts: make(map[string]int),
	}
}

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.0001
}

func TestFinalizeSourceStats_EvenSplit(t *testing.T) {
	b := newTestBuilder([]domain.PlaylistSource{
		{Name: "cbs", ID: "A"},
		{Name: "nbc", ID: "B"},
	})
	b.recordVideoSource("v1", "A")
	b.recordVideoSource("v2", "B")
	b.recordSourceContribution("dist/v1.mp4", 30*time.Second, 3)
	b.recordSourceContribution("dist/v2.mp4", 30*time.Second, 3)

	m := &domain.Manifest{}
	b.finalizeSourceStats(m)

	if len(m.Sources) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(m.Sources))
	}
	for _, s := range m.Sources {
		if !approxEqual(s.Percentage, 50) {
			t.Errorf("source %s: expected 50%%, got %v", s.Name, s.Percentage)
		}
		if s.ClipCount != 3 {
			t.Errorf("source %s: expected 3 clips, got %d", s.Name, s.ClipCount)
		}
		if !approxEqual(s.DurationSeconds, 30) {
			t.Errorf("source %s: expected 30s, got %v", s.Name, s.DurationSeconds)
		}
	}
}

func TestFinalizeSourceStats_PreservesConfigOrder(t *testing.T) {
	b := newTestBuilder([]domain.PlaylistSource{
		{Name: "cbs", ID: "A"},
		{Name: "nbc", ID: "B"},
		{Name: "abc", ID: "C"},
	})
	b.recordVideoSource("v1", "C")
	b.recordSourceContribution("v1.mp4", 10*time.Second, 1)

	m := &domain.Manifest{}
	b.finalizeSourceStats(m)

	want := []string{"cbs", "nbc", "abc"}
	for i, s := range m.Sources {
		if s.Name != want[i] {
			t.Errorf("position %d: expected %q, got %q", i, want[i], s.Name)
		}
	}
}

func TestFinalizeSourceStats_ZeroContentSource(t *testing.T) {
	b := newTestBuilder([]domain.PlaylistSource{
		{Name: "cbs", ID: "A"},
		{Name: "nbc", ID: "B"},
	})
	b.recordVideoSource("v1", "A")
	b.recordSourceContribution("v1.mp4", 60*time.Second, 4)

	m := &domain.Manifest{}
	b.finalizeSourceStats(m)

	byID := map[string]domain.SourceStat{}
	for _, s := range m.Sources {
		byID[s.PlaylistID] = s
	}
	if !approxEqual(byID["A"].Percentage, 100) {
		t.Errorf("source A: expected 100%%, got %v", byID["A"].Percentage)
	}
	if !approxEqual(byID["B"].Percentage, 0) {
		t.Errorf("source B: expected 0%%, got %v", byID["B"].Percentage)
	}
	if byID["B"].ClipCount != 0 {
		t.Errorf("source B: expected 0 clips, got %d", byID["B"].ClipCount)
	}
}

func TestFinalizeSourceStats_NoContent(t *testing.T) {
	b := newTestBuilder([]domain.PlaylistSource{
		{Name: "cbs", ID: "A"},
	})
	m := &domain.Manifest{}
	b.finalizeSourceStats(m)

	if len(m.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(m.Sources))
	}
	if m.Sources[0].Percentage != 0 {
		t.Errorf("expected 0%% (no division by zero), got %v", m.Sources[0].Percentage)
	}
}

func TestRecordSourceContribution_UnknownVideoLogsButNoCrash(t *testing.T) {
	b := newTestBuilder([]domain.PlaylistSource{{Name: "cbs", ID: "A"}})
	b.recordSourceContribution("unknown.mp4", 10*time.Second, 1)
	if len(b.sourceDurations) != 0 {
		t.Errorf("expected no durations recorded for unknown video, got %v", b.sourceDurations)
	}
}
