package downloader

import (
	"context"
	"errors"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/youtubeapi"
)

func (d *Downloader) GetChannelPlaylistID(ctx context.Context, userName string) (id string, err error) {

	resp, err := d.ytapi.ChannelsList(ctx, youtubeapi.ChannelsListRequest{
		ForHandle: userName,
		Part:      []string{"contentDetails"},
	})

	if len(resp.Items) == 0 {
		err = errors.New("no channels in response")
		log.Error(err.Error())
		return
	}

	id = resp.Items[0].ContentDetails.RelatedPlaylists.Uploads
	if id == "" {
		err = errors.New("recieved empty id")
		log.Error(err.Error())
	}

	return
}
