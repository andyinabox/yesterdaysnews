package errorhandler

import (
	"context"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func DefaultErrorHandler(ctx context.Context, totalDownloadTargets int) domain.ErrorHandler {

	fatalFunc := func(typ string, err error) {
		log.Fatalf("%s: %s", typ, err)
	}

	errorFunc := func(typ string, err error) {
		log.Errorf("%s: %s", typ, err)
	}

	return New(ctx, &Config{
		ErrorFunc: errorFunc,
		FatalFunc: fatalFunc,
		Thresholds: map[string]int{
			domain.ErrTypeFatal:         0,
			domain.ErrTypeMoveObject:    totalDownloadTargets / 3,
			domain.ErrTypeDownloadVideo: totalDownloadTargets / 3,
			domain.ErrTypeCutVideo:      totalDownloadTargets / 3,
		},
	})
}
