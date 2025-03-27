package builder

import (
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func (b *Builder) createManifest(buildId string) *domain.Manifest {
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
