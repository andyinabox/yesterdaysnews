package builder

import (
	"strings"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

// ParsePlaylists parses a comma-separated list of playlist sources. Each entry
// may be either a bare playlist ID or a "name:id" pair. Bare IDs use the ID as
// the name. Empty entries and entries with empty IDs are skipped.
func ParsePlaylists(raw string) []domain.PlaylistSource {
	parts := strings.Split(raw, ",")
	sources := make([]domain.PlaylistSource, 0, len(parts))
	for _, part := range parts {
		entry := strings.TrimSpace(part)
		if entry == "" {
			continue
		}
		name, id := entry, entry
		if idx := strings.Index(entry, ":"); idx >= 0 {
			name = strings.TrimSpace(entry[:idx])
			id = strings.TrimSpace(entry[idx+1:])
			if name == "" {
				name = id
			}
		}
		if id == "" {
			continue
		}
		sources = append(sources, domain.PlaylistSource{Name: name, ID: id})
	}
	return sources
}
