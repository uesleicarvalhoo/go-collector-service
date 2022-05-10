package sender

import (
	"context"
	"io"

	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
)

type Streamer interface {
	NotifyPublishedFile(ctx context.Context, fileKey string, file models.File) error
	NotifyInvalidFile(ctx context.Context, file models.File) error
}

type Storage interface {
	SendFile(ctx context.Context, fileKey string, reader io.ReadSeeker) (err error)
}

type FileServer interface {
	Glob(ctx context.Context, pattern string) ([]string, error)
	Open(ctx context.Context, filePath string) (io.ReadSeekCloser, error)
	Remove(ctx context.Context, filePath string) error
}
