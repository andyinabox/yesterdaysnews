package response

// type ResourceType string

// const (
// 	ResourceTypePlaylistItem ResourceType = "PlaylistItem"
// 	ResourceTypeChannel      ResourceType = "Channel"
// 	ResourceTypeVideo        ResourceType = "Video"
// )

type Response struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	PrevPageToken string `json:"prevPageToken"`
	// Items         json.RawMessage `json:"items"`
	PageInfo PageInfo `json:"pageInfo"`
}

// var validKinds map[ResourceType][]string

// func init() {
// 	validKinds = map[ResourceType][]string{}
// 	validKinds[ResourceTypePlaylistItem] = []string{"youtube#playlistItemListResponse"}
// 	validKinds[ResourceTypeChannel] = []string{"youtube#channelListResponse"}
// 	validKinds[ResourceTypeVideo] = []string{"youtube#videoListResponse"}
// }

// func isValidKind(typ ResourceType, kind string) error {
// 	kinds, ok := validKinds[typ]
// 	if !ok {
// 		return fmt.Errorf("invalid type: %q", typ)
// 	}

// 	for _, k := range kinds {
// 		if k == kind {
// 			return nil
// 		}
// 	}

// 	return fmt.Errorf("response of kind %q does not return type %s", kind, typ)
// }

// func (s *Success) PlaylistItems() (items []PlaylistItem, err error) {

// 	err = isValidKind(ResourceTypePlaylistItem, s.Kind)
// 	if err != nil {
// 		return
// 	}

// 	items = make([]PlaylistItem, 0)
// 	err = json.Unmarshal(s.Items, &items)
// 	return
// }

// func (s *Success) ChannelItems() (items []Channel, err error) {

// 	err = isValidKind(ResourceTypeChannel, s.Kind)
// 	if err != nil {
// 		return
// 	}

// 	items = make([]Channel, 0)
// 	err = json.Unmarshal(s.Items, &items)
// 	return
// }

// func (s Success) VideoItems() (items []Video, err error) {
// 	err = isValidKind(ResourceTypeVideo, s.Kind)
// 	if err != nil {
// 		return
// 	}

// 	items = make([]Video, 0)
// 	err = json.Unmarshal(s.Items, &items)
// 	return
// }
