package youtubeapi

import (
	"encoding/json"
	"time"
)

type successResponse struct {
	Kind          string          `json:"kind"`
	Etag          string          `json:"etag"`
	NextPageToken string          `json:"nextPageToken"`
	PrevPageToken string          `json:"prevPageToken"`
	Items         json.RawMessage `json:"items"`
	PageInfo      pageInfo        `json:"pageInfo"`
}

// generic types

type thumbnailData struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type resourceID struct {
	VideoID string `json:"videoId"`
}

type pageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

// playlistId response

type playlistIdItem struct {
	Kind    string            `json:"kind"`
	Snippet playlistIdSnippet `json:"snippet"`
}

type playlistIdSnippet struct {
	PublishedAt time.Time                `json:"publishedAt"`
	ResourceID  resourceID               `json:"resourceId"`
	Thumbnails  map[string]thumbnailData `json:"thumbnails"`
}

type playlistIdContentDetails struct {
	VideoId          string    `json:"videoId"`
	VideoPublishedAt time.Time `json:"videoPublishedAt"`
}
