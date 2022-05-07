package storage

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
)

type NoneStorage struct{}

func (ns *NoneStorage) SendFile(ctx context.Context, fileKey string, reader io.ReadSeeker) error {
	logrus.Infof("File received, %+v", fileKey)

	return nil
}

func NewNoneStorage() *NoneStorage {
	return &NoneStorage{}
}
