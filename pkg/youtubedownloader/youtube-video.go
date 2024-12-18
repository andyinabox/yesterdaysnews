package youtubedownloader

type YouTubeVideo struct {
	// this is only a subset of fields I thought might be useful
	// see esample in test/video-info.json
	ID                 string                     `json:"id"`
	Title              string                     `json:"title"`
	Format             string                     `json:"format"`
	Width              uint                       `json:"width"`
	Height             uint                       `json:"height"`
	AspectRatio        float32                    `json:"aspect_ratio"`
	Ext                string                     `json:"ext"`
	Filename           string                     `json:"filename"`
	UploadDate         string                     `json:"upload_date"` // format: YYYMMDD
	Timestamp          uint                       `json:"timestamp"`   // format: 1734177608 (js timestamp i think)
	Formats            []YouTubeVideoFormat       `json:"formats"`
	AutomaticCaptions  map[string]YouTubeSubtitle `json:"automatic_captions"`
	RequestedSubtitles map[string]YouTubeSubtitle `json:"requested_formats"`
}

type YouTubeVideoFormat struct {
	Format string `json:"format"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
	Ext    string `json:"ext"`
	URL    string `json:"url"`
}

type YouTubeSubtitle struct {
	URL      string `json:"url"`
	Ext      string `json:"ext"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
}
