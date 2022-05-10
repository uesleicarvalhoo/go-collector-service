package sender

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
)

var (
	ErrInvalidWorkersCount  = errors.New("workers must be higher then 0")
	ErrInvalidPattern       = errors.New("one or more patterns must be informed")
	ErrServiceAlreadStarted = errors.New("service is running")
)

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

type Config struct {
	EventTopic    string
	MatchPatterns []string
	Workers       int
	Delay         time.Duration
}
