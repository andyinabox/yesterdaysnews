package youtubeservice

import (
	"context"

	"code.andydayton.com/andy/yesterdaysnews/pkg/youtubeapi"
)

func (s *Service) GetChannelPlaylistID(ctx context.Context, userName string) (id string, err error) {

	resp, err := s.ytapi.ChannelsList(ctx, youtubeapi.ChannelsListRequest{
		ForHandle: userName,
		Part:      []string{"contentDetails"},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Items) == 0 {
		return "", err
	}

	id = resp.Items[0].ContentDetails.RelatedPlaylists.Uploads
	if id == "" {
		return "", err
	}

	return
}
