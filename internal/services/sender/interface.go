package sender

import (
	"context"
	"io"

	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
)

type Streamer interface {
	NotifyPublishedFile(fileKey string, file models.File) error
	NotifyInvalidFile(file models.File) error
}

type Storage interface {
	SendFile(ctx context.Context, fileKey string, reader io.ReadSeeker) (err error)
}

type Collector interface {
	// GetNextFile() (models.File, error)
	GetFiles() ([]models.File, error)
	RemoveFile(file models.File) error
}
