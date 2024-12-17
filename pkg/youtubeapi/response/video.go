package response

type Video struct {
	Kind string `json:"kind"`
	Etag string `json:"etag"`
	ID   string `json:"id"`
}

type VideoFileDetails struct {
	VideoStreams []VideoStream `json:"videoStreams"`
}

type VideoStream struct {
	AspectRatio  float32 `json:"aspectRatio"`
	WidthPixels  uint    `json:"widthPixels"`
	HeightPixels uint    `json:"heightPixels"`
}
