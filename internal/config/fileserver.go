package config

type FileServerConfig struct {
	Server     string `envconfig:"FILE_SERVER_URL" default:"localhost:22"`
	User       string `envconfig:"FILE_SERVER_USER" default:"admin"`
	Password   string `envconfig:"FILE_SERVER_PASSWORD" default:"secret"`
	PrivateKey string `envconfig:"FILE_SERVER_PRIVATE_KEY"`

	KeyExchanges string `envconfig:"FILE_SERVER_KEY_EXCHANGES"`
}
