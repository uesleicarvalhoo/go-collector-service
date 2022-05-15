package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	ServiceName    string `envconfig:"SERVICE_NAME" default:"go-collector"`
	ServiceVersion string `envconfig:"SERVICE_VERSION" default:"0.0.0"`
	Env            string `envconfig:"ENVIRONMENT" default:"dev"`
	Debug          bool   `envconfig:"DEBUG" default:"false"`

	AwsRegion string `envconfig:"AWS_REGION" default:"sa-east-1"`

	TraceServiceName string `envconfig:"TRACE_SERVICE_NAME"`
	TraceURL         string `envconfig:"TRACE_URL" default:"http://localhost:14268"`
	TraceEnable      bool   `envconfig:"TRACE_ENABLE" defult:"false"`

	BrokerConfig     BrokerConfig
	StorageConfig    StorageConfig
	FileServerConfig FileServerConfig
	LoggerConfig     LoggerConfig
}

func (s *Settings) LoadFromEnv() error {
	err := godotenv.Load(os.Getenv("ENVFILE_PATH"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = envconfig.Process("", s)
	if err != nil {
		return err
	}

	if strings.TrimSpace(s.TraceServiceName) == "" {
		s.TraceServiceName = s.ServiceName
	}

	return nil
}
