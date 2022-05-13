package sender

import (
	"github.com/uesleicarvalhoo/go-collector-service/internal/config"
	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
)

type Config = config.SenderConfig

func validateConfig(cfg Config) error {
	validator := models.Validator{}

	if cfg.ParalelUploads == 0 {
		validator.AddError(models.ValidationErrorProps{Context: "config", Message: "'ParalelUploads' must be higher then 0"})
	}

	if len(cfg.MatchPatterns) == 0 {
		validator.AddError(
			models.ValidationErrorProps{Context: "config", Message: "'MatchPatterns' must be have one or more patterns"},
		)
	}

	if validator.HasErrors() {
		return validator.GetError()
	}

	return nil
}
