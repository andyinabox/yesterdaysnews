package builder

import (
	"time"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/util"
)

func (b *Builder) createManifest() *domain.Manifest {
	now := time.Now()
	yesterday := util.Yesterday()
	return &domain.Manifest{
		BuildDate:   now,
		ContentDate: yesterday,
		ID:          util.Timestamp(now),
		Files: domain.ManifestFiles{
			Clips: []string{},
		},
	}
}
