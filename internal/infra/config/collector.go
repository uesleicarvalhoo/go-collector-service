package config

type CollectorConfig struct {
	MatchPattern string `env:"COLLECTOR_MATCH_PATTERN,default=/files/*"`
	Server       string `env:"COLLECTOR_SERVER_URL,default=localhost:22"`
	User         string `env:"COLLECTOR_SERVER_USER,default=admin"`
	Password     string `env:"COLLECTOR_SERVER_PASSWORD,default=secret"`
	PrivateKey   string `env:"COLLECTOR_SERVER_PRIVATE_KEY"`

	KeyExchanges string `env:"COLLECTOR_KEY_EXCHANGES"`
}
