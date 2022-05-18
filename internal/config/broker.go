package config

type BrokerConfig struct {
	Host     string `envconfig:"BROKER_URL" default:"localhost"`
	Port     string `envconfig:"BROKER_PORT" default:"5672"`
	User     string `envconfig:"BROKER_USER" default:"guest"`
	Password string `envconfig:"BROKER_PASSWORD" default:"guest"`
}
