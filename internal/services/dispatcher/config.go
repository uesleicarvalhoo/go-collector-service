package dispatcher

import (
	"fmt"
	"io/ioutil"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/sender"
	"gopkg.in/yaml.v3"
)

type Config struct {
	SenderConfig []sender.Config `yaml:"sender" json:"sender"`
}

func (c *Config) LoadFromYaml(configpath string) error {
	data, err := ioutil.ReadFile(configpath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return c.Validate()
}

func (c Config) Validate() error {
	validator := models.Validator{}

	if len(c.SenderConfig) == 0 {
		validator.AddError("SenderConfig", "sender config is required")
	}

	for nWorker, cfg := range c.SenderConfig {
		if err := cfg.Validate(); err != nil {
			validator.AddError(fmt.Sprintf("Worker[%d]", nWorker+1), err.Error())
		}
	}

	if validator.HasErrors() {
		return validator.GetError()
	}

	return nil
}
