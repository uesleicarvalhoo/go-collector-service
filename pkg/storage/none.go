package storage

import (
	"context"
	"io"

	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
)

type NoneStorage struct{}

func (ns *NoneStorage) SendFile(ctx context.Context, fileKey string, reader io.ReadSeeker) error {
	logger.Debugf("File received, %+v", fileKey)

	return nil
}

func NewNoneStorage() *NoneStorage {
	return &NoneStorage{}
}
