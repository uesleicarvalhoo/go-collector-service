package config

type FileServerConfig struct {
	Server     string `env:"FILE_SERVER_URL,default=localhost:22"`
	User       string `env:"FILE_SERVER_USER,default=admin"`
	Password   string `env:"FILE_SERVER_PASSWORD,default=secret"`
	PrivateKey string `env:"FILE_SERVER_PRIVATE_KEY"`

	KeyExchanges string `env:"FILE_SERVER_KEY_EXCHANGES"`
}
