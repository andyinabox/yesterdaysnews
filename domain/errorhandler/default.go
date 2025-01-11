package errorhandler

import (
	"context"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func DefaultErrorHandler(ctx context.Context) domain.ErrorHandler {

	errorFunc := func(typ string, err error) {
		if typ == domain.ErrTypeFatal {
			log.Fatalf("%s: %s", typ, err)
		}
		log.Errorf("%s: %s", typ, err)
	}

	return New(ctx, &Config{
		ErrorFunc: errorFunc,
		Thresholds: map[string]int{
			domain.ErrTypeMoveObject: 10,
		},
	})
}
