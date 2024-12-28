package osidentityclient

type Config struct {
	IdentityEndpoint string
	UserDomainName   string
	Username         string
	Password         string
}

type Client struct {
	cfg *Config
}

func New(cfg *Config) *Client {
	return &Client{cfg}
}
