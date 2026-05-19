package builder

import (
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (b *Builder) recordVideoSource(videoID, playlistID string) {
	b.sourceMu.Lock()
	defer b.sourceMu.Unlock()
	b.videoSourceMap[videoID] = playlistID
}

func (b *Builder) recordFilterReason(playlistID string, reason domain.FilterReason) {
	b.sourceMu.Lock()
	defer b.sourceMu.Unlock()
	stats, ok := b.sourceFilterStats[playlistID]
	if !ok {
		stats = make(map[domain.FilterReason]int)
		b.sourceFilterStats[playlistID] = stats
	}
	stats[reason]++
}

func (b *Builder) recordSourceContribution(videoFilePath string, duration time.Duration, clipCount int) {
	base := filepath.Base(videoFilePath)
	videoID := strings.TrimSuffix(base, filepath.Ext(base))

	b.sourceMu.Lock()
	defer b.sourceMu.Unlock()

	playlistID, ok := b.videoSourceMap[videoID]
	if !ok {
		slog.Warn("no source recorded for video", "videoID", videoID, "file", videoFilePath)
		return
	}
	b.sourceDurations[playlistID] += duration
	b.sourceClipCounts[playlistID] += clipCount
}

// finalizeSourceStats walks the configured playlists in order and writes
// per-source totals onto the manifest. It also emits a structured log line per
// source so the breakdown is visible in build output.
func (b *Builder) finalizeSourceStats(manifest *domain.Manifest) {
	b.sourceMu.Lock()
	defer b.sourceMu.Unlock()

	var total time.Duration
	for _, d := range b.sourceDurations {
		total += d
	}

	stats := make([]domain.SourceStat, 0, len(b.cfg.Playlists))
	for _, p := range b.cfg.Playlists {
		dur := b.sourceDurations[p.ID]
		var pct float64
		if total > 0 {
			pct = float64(dur) / float64(total) * 100
		}
		stat := domain.SourceStat{
			Name:            p.Name,
			PlaylistID:      p.ID,
			DurationSeconds: dur.Seconds(),
			Percentage:      pct,
			ClipCount:       b.sourceClipCounts[p.ID],
		}
		stats = append(stats, stat)
		slog.Info("source breakdown",
			"name", stat.Name,
			"playlistID", stat.PlaylistID,
			"percent", stat.Percentage,
			"durationSeconds", stat.DurationSeconds,
			"clips", stat.ClipCount,
		)
	}
	manifest.Sources = stats

	for _, p := range b.cfg.Playlists {
		fs := b.sourceFilterStats[p.ID]
		slog.Info("filter breakdown",
			"name", p.Name,
			"playlistID", p.ID,
			"fetched", fs[domain.FilterReasonFetched],
			"wrongDate", fs[domain.FilterReasonWrongDate],
			"formatUnavailable", fs[domain.FilterReasonFormatUnavailable],
			"tooLarge", fs[domain.FilterReasonTooLarge],
			"noCaptions", fs[domain.FilterReasonNoCaptions],
			"infoError", fs[domain.FilterReasonInfoError],
			"accepted", fs[domain.FilterReasonAccepted],
		)
	}
}
