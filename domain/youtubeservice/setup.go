package youtubeservice

import "context"

func (s *Service) Setup(ctx context.Context) error {
	return s.ytdl.CheckVersion(ctx)
}
