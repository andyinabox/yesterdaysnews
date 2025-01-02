package clipstreamer

import (
	"context"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func (s *Streamer) ErrorStream(ctx context.Context, handler func(domain.StreamErr)) chan<- domain.StreamErr {
	stream := make(chan domain.StreamErr)

	go func() {
		defer close(stream)
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-stream:
				handler(err)
			}
		}
	}()

	return stream
}
