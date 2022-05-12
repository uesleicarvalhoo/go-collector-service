package sender

import (
	"context"
	"errors"
	"io"

	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
	"github.com/uesleicarvalhoo/go-collector-service/internal/infra/config"
)

var ErrServiceAlreadStarted = errors.New("service is running")

type Config = config.SenderConfig

func validateConfig(cfg Config) error {
	validator := models.Validator{}

	if cfg.MaxCollectBatchSize == 0 {
		validator.AddError(
			models.ValidationErrorProps{Context: "config", Message: "'maxCollectBatchSize' must be higher then 0"},
		)
	}

	if cfg.ParalelUploads == 0 {
		validator.AddError(models.ValidationErrorProps{Context: "config", Message: "'PralelUploads' must be higher then 0"})
	}

	if len(cfg.MatchPatterns) == 0 {
		validator.AddError(
			models.ValidationErrorProps{Context: "config", Message: "'MatchPatterns' must be have one or more pattern"},
		)
	}

	if validator.HasErrors() {
		return validator.GetError()
	}

	return nil
}

type Storage interface {
	SendFile(context.Context, string, io.ReadSeeker) (err error)
}

type FileServer interface {
	Glob(context.Context, string) ([]string, error)
	Open(context.Context, string) (io.ReadSeekCloser, error)
	MoveFile(context.Context, string, string) error
}

type Broker interface {
	SendEvent(models.Event) error
}
