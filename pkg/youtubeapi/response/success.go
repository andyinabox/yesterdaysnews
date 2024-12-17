package response

import (
	"encoding/json"
)

type Success struct {
	Kind          string          `json:"kind"`
	Etag          string          `json:"etag"`
	NextPageToken string          `json:"nextPageToken"`
	PrevPageToken string          `json:"prevPageToken"`
	Items         json.RawMessage `json:"items"`
	PageInfo      PageInfo        `json:"pageInfo"`
}

func (s *Success) PlaylistItems() (items []PlaylistItem, err error) {
	items = make([]PlaylistItem, 0)
	err = json.Unmarshal(s.Items, &items)
	return
}

func (s *Success) ChannelItems() (items []Channel, err error) {
	items = make([]Channel, 0)
	err = json.Unmarshal(s.Items, &items)
	return
}
