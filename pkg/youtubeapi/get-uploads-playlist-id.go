package youtubeapi

import (
	"context"
	"errors"
	"net/url"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/youtubeapi/response"
)

func (c *Client) GetUploadsPlaylistIdForUser(ctx context.Context, userName string) (id string, err error) {

	q := make(url.Values)
	q.Add("part", "contentDetails")
	q.Add("forHandle", userName)

	var resp *response.Success
	resp, err = c.doGetRequest(ctx, "channels", q)
	if err != nil {
		return
	}

	var channels []response.Channel
	channels, err = resp.ChannelItems()
	if err != nil {
		return
	}

	if len(channels) == 0 {
		err = errors.New("no channels in response")
		log.Error(err.Error())
		return
	}

	id = channels[0].ContentDetails.RelatedPlaylists.Uploads
	if id == "" {
		err = errors.New("recieved empty id")
		log.Error(err.Error())
	}

	return
}
