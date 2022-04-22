package services

import (
	"context"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
)

type Streamer interface {
	NotifyPublishedFile(file models.File) error
	NotifyInvalidFile(file models.File) error
}

type Storage interface {
	SendFile(ctx context.Context, file models.File) (err error)
}

type Collector interface {
	ListFiles() ([]models.FileInfo, error)
	GetFileData(file models.FileInfo) ([]byte, error)
	RemoveFile(file models.FileInfo) error
}
