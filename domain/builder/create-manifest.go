package builder

import (
	"log/slog"
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func (b *Builder) createManifest(buildId string) *domain.Manifest {
	slog.Debug("create manifest", "id", buildId)
	now := time.Now()
	if buildId == "" {
		buildId = util.Timestamp(now)
	}
	yesterday := util.Yesterday()
	return &domain.Manifest{
		BuildDate:   now,
		ContentDate: yesterday,
		ID:          buildId,
		Files: domain.ManifestFiles{
			Clips: []string{},
		},
	}
}
