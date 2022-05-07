package config

type BrokerConfig struct {
	Host       string `env:"BROKER_URL,default=amqp://localhost"`
	Port       string `env:"BROKER_PORT,default=5672"`
	User       string `env:"BROKER_USER,default=guest"`
	Password   string `env:"BROKER_PASSWORD,default=guest"`
	EventTopic string `env:"BROKER_EVENT_TOPIC,default=collector.files"`
}
