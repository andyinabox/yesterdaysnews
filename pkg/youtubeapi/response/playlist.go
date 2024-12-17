package response

import "time"

type PlaylistItem struct {
	Kind           string                     `json:"kind"`
	Etag           string                     `json:"etag"`
	ID             string                     `json:"id"`
	Snippet        PlaylistItemSnippet        `json:"snippet"`
	ContentDetails PlaylistItemContentDetails `json:"contentDetails"`
	Status         PlaylistItemStatus         `json:"status"`
}

type PlaylistItemSnippet struct {
	PublishedAt time.Time                `json:"publishedAt"`
	Thumbnails  map[string]ThumbnailData `json:"thumbnails"`
	// many others that haven't been implemented
	ResourceID ResourceID `json:"resourceId"`
}

type PlaylistItemContentDetails struct {
	VideoId          string    `json:"videoId"`
	VideoPublishedAt time.Time `json:"videoPublishedAt"`
}

type PlaylistItemStatus struct {
	PrivacyStatus string `json:"privacyStatus"`
}
