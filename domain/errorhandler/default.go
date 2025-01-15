package errorhandler

import (
	"context"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func DefaultErrorHandler(ctx context.Context, totalDownloadTargets int) domain.ErrorHandler {
	return New(ctx, &Config{
		Thresholds: map[string]int{
			domain.ErrTypeFatal:         0,
			domain.ErrTypeMoveObject:    totalDownloadTargets / 3,
			domain.ErrTypeDownloadVideo: totalDownloadTargets / 3,
			domain.ErrTypeCutVideo:      totalDownloadTargets / 3,
		},
	})
}
