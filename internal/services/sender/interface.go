package sender

import (
	"context"
	"errors"
	"io"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/fileserver"
)

type LockInterface = fileserver.LockerInterface

var ErrServiceAlreadStarted = errors.New("service is running")

type Storage interface {
	SendFile(context.Context, string, io.ReadSeeker) (err error)
}

type FileServer interface {
	Glob(context.Context, string) ([]string, error)
	Open(context.Context, string) (io.ReadSeekCloser, error)
	MoveFile(context.Context, string, string) error
	Lock(context.Context, string) (LockInterface, error)
}

type Broker interface {
	SendEvent(models.Event) error
}
