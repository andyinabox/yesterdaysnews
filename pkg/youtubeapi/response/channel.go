package response

type ChannelListResponse struct {
	Response
	Items []Channel `json:"items"`
}

type Channel struct {
	Kind           string                `json:"kind"`
	Etag           string                `json:"etag"`
	ID             string                `json:"id"`
	ContentDetails ChannelContentDetails `json:"contentDetails"`
	// several more that I'm not using yet
}

type ChannelContentDetails struct {
	RelatedPlaylists ChannelRelatedPlaylists `json:"relatedPlaylists"`
}

type ChannelRelatedPlaylists struct {
	Likes     string `json:"likes"`
	Favorites string `json:"favorites"`
	Uploads   string `json:"uploads"`
}
