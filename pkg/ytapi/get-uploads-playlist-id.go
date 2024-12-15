package ytapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
)

type relatedPlaylists struct {
	Uploads string `json:"uploads"`
}

type contentDetails struct {
	RelatedPlaylists relatedPlaylists `json:"relatedPlaylists"`
}

type getUploadsPlaylistIdForUserRespItem struct {
	Kind           string         `json:"kind"`
	ContentDetails contentDetails `json:"contentDetails"`
}

type getUploadsPlaylistIdForUserResp struct {
	Kind  string                                `json:"kind"`
	Items []getUploadsPlaylistIdForUserRespItem `json:"items"`
}

func (c *Client) GetUploadsPlaylistIdForUser(ctx context.Context, userName string) (id string, err error) {

	q := make(url.Values)
	q.Add("part", "contentDetails")
	q.Add("forHandle", userName)

	var resp *http.Response
	resp, err = c.doGetRequest(ctx, "channels", q)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Debug(string(body))

	data := getUploadsPlaylistIdForUserResp{}
	err = json.Unmarshal(body, &data)

	if len(data.Items) == 0 {
		err = errors.New("no items in response")
		log.Error(err.Error(), "body", data)
		return
	}

	id = data.Items[0].ContentDetails.RelatedPlaylists.Uploads
	if id == "" {
		err = errors.New("recieved empty id")
		log.Error(err.Error(), "data", data)
	}

	return
}
