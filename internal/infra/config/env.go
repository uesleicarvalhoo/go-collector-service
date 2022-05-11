package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
)

type AppSettings struct {
	ServiceName    string `envconfig:"SERVICE_NAME" default:"go-collector"`
	ServiceVersion string `envconfig:"SERVICE_VERSION" default:"0.0.0"`
	Env            string `envconfig:"ENVIRONMENT" default:"dev"`
	Debug          bool   `envconfig:"DEBUG" default:"false"`

	AwsRegion string `envconfig:"AWS_REGION" default:"sa-east-1"`

	TraceServiceName string `envconfig:"TRACE_SERVICE_NAME"`
	TraceURL         string `envconfig:"TRACE_URL" default:"http://localhost:14268"`

	EventTopic     string   `envconfig:"SENDER_EVENT_TOPIC" default:"collector.services"`
	ParalelUploads int      `envconfig:"SENDER_PARALLEL_UPLOADS" default:"2"`
	CollectDelay   int      `envconfig:"SENDER_COLLECT_DELAY" default:"5"`
	MatchPatterns  []string `envconfig:"SENDER_MATCH_PATTERNS" required:"true" default:"upload/*.json"`

	BrokerConfig     BrokerConfig
	StorageConfig    StorageConfig
	FileServerConfig FileServerConfig
}

func LoadAppSettingsFromEnv() AppSettings {
	var cfg AppSettings

	err := godotenv.Load()
	if err != nil {
		logger.Infof("Couldn't be load env from .env file: %s", err)
	}

	err = envconfig.Process("", &cfg)
	if err != nil {
		logger.Fatal(err)

		return AppSettings{}
	}

	if strings.TrimSpace(cfg.TraceServiceName) == "" {
		cfg.TraceServiceName = cfg.ServiceName
	}

	return cfg
}
