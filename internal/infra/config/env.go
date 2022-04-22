package config

import (
	"log"

	"github.com/netflix/go-env"
)

type AppSettings struct {
	Env   string `env:"ENVIRONMENT,default=dev"`
	Debug bool   `env:"DEBUG,default=false"`

	AwsRegion string `env:"AWS_REGION,default=sa-east-1"`

	TraceServiceName string `env:"TRACE_SERVICE_NAME"`
	TraceURL         string `env:"TRACE_URL,default=http://localhost:14268"`

	BrokerConfig  BrokerConfig
	StorageConfig StorageConfig
}

func LoadAppSettingsFromEnv() (cfg AppSettings) {
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		log.Fatal(err)

		return
	}

	if cfg.TraceServiceName == "" {
		cfg.TraceServiceName = ServiceName
	}

	return
}
