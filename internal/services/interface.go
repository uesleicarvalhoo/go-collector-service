package services

import (
	"context"
	"io"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
)

type FileServer interface {
	Glob(context.Context, string) ([]string, error)
	Open(context.Context, string) (io.ReadSeekCloser, error)
	Move(context.Context, string, string) error
	AcquireLock(context.Context, string) error
	ReleaseLock(context.Context, string) error
}
type Storage interface {
	SendFile(context.Context, string, io.ReadSeeker) (err error)
}

type Broker interface {
	SendEvent(models.Event) error
}
