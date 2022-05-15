package collector

import (
	"fmt"
	"strings"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
)

type Config struct {
	// Patterns to collect files on file server
	// Ex: ./files/*.json"
	MatchPatterns []string `yaml:"pattern" json:"pattern"`
	// Max files amount to collect on each collect loop
	MaxCollectBatchSize int `yaml:"maxFilesBatch" json:"maxFilesBatch"`
	// Minimum seconds between each collect loop
	CollectDelay int `yaml:"delay" json:"delay"`
}

func (c *Config) Validate() error {
	validator := models.Validator{}

	if len(c.MatchPatterns) == 0 {
		validator.AddError("MatchPatterns", "field is required")
	}

	for i, pattern := range c.MatchPatterns {
		if strings.TrimSpace(pattern) == "" {
			validator.AddError(fmt.Sprintf("MatchPattern[%d]", i), "field is required")
		}
	}

	if validator.HasErrors() {
		return validator.GetError()
	}

	return nil
}
