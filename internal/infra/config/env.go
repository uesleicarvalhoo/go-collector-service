package config

import (
	"github.com/joho/godotenv"
	"github.com/netflix/go-env"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
)

type AppSettings struct {
	Env   string `env:"ENVIRONMENT,default=dev"`
	Debug bool   `env:"DEBUG,default=false"`

	AwsRegion string `env:"AWS_REGION,default=sa-east-1"`

	TraceServiceName string `env:"TRACE_SERVICE_NAME"`
	TraceURL         string `env:"TRACE_URL,default=http://localhost:14268"`

	MatchPattern string `env:"SENDER_MATCH_PATTERN,default=/files/*"`

	BrokerConfig     BrokerConfig
	StorageConfig    StorageConfig
	FileServerConfig FileServerConfig
}

func LoadAppSettingsFromEnv() AppSettings {
	var cfg AppSettings

	err := godotenv.Load()
	if err != nil {
		logger.Info("Couldn't be load env from .env file")
	}

	_, err = env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		logger.Fatal(err)

		return AppSettings{}
	}

	if cfg.TraceServiceName == "" {
		cfg.TraceServiceName = ServiceName
	}

	return cfg
}
