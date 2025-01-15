package identity

type Config struct {
	IdentityEndpoint string
	UserID           string
	Password         string
}

type Client struct {
	cfg *Config
}

func New(cfg *Config) *Client {
	return &Client{cfg}
}
