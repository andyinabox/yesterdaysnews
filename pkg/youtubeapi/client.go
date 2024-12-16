package youtubeapi

const APIBase = "https://youtube.googleapis.com/youtube/v3"

type Client struct {
	key string
}

func New(apiKey string) *Client {
	return &Client{
		key: apiKey,
	}
}
