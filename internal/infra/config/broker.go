package config

type BrokerConfig struct {
	Host     string `env:"BROKER_URL,default=localhost"`
	Port     string `env:"BROKER_PORT,default=5672"`
	User     string `env:"BROKER_USER,default=guest"`
	Password string `env:"BROKER_PASSWORD,default=guest"`
}
