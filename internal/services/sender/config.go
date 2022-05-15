package sender

import (
	"strings"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/collector"
)

type Config struct {
	// Broker topic name to send event with file process result
	EventTopic string `yaml:"topic" json:"topic"`

	// Number of workers to publish file to storage
	Workers      int              `yaml:"workers" json:"workers"`
	CollectorCfg collector.Config `json:"collect" yaml:"collect"`
}

func (c Config) Validate() error {
	validator := models.Validator{}

	if strings.TrimSpace(c.EventTopic) == "" {
		validator.AddError("eventTopic", "field is required")
	}

	if c.Workers == 0 {
		validator.AddError("workers", "must be higher then 0")
	}

	if err := c.CollectorCfg.Validate(); err != nil {
		validator.AddError("collector", err.Error())
	}

	if validator.HasErrors() {
		return validator.GetError()
	}

	return nil
}
