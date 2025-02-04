package builder

import (
	"context"
	"time"
)

func (b *Builder) ExtractImages(ctx context.Context, clips <-chan string) <-chan string {
	return b.vp.ExtractImagesStream(ctx, b.errs, clips, time.Duration(0))
}
