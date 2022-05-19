package services

import (
	"context"
	"io"
	"io/fs"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
)

type Locker = models.Locker

type FileServer interface {
	Glob(context.Context, string) ([]string, error)
	Open(context.Context, string) (io.ReadSeekCloser, error)
	Move(context.Context, string, string) error
	Stat(context.Context, string) (fs.FileInfo, error)
	AcquireLock(context.Context, string) (Locker, error)
}

type Storage interface {
	SendFile(context.Context, string, io.ReadSeeker) (err error)
}

type Broker interface {
	SendEvent(models.Event) error
}
