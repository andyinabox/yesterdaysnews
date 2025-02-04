package errorhandler

import (
	"context"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func DefaultErrorHandler(ctx context.Context, totalDownloadTargets int) domain.ErrorHandler {
	return New(ctx, &Config{
		Thresholds: map[string]int{
			domain.ErrTypeFatal:         0,
			domain.ErrTypeTODO:          0,
			domain.ErrTypeMoveObject:    totalDownloadTargets / 3,
			domain.ErrTypeDownloadVideo: totalDownloadTargets / 3,
			domain.ErrTypeGetVideoID:    totalDownloadTargets / 3,
			domain.ErrTypeUploadFile:    totalDownloadTargets / 3,
			domain.ErrTypeCutVideo:      50,
			domain.ErrTypeAverageImage:  50,
		},
	})
}
