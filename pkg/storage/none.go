package storage

import (
	"context"
	"log"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
)

type NoneStorage struct{}

func (ns *NoneStorage) SendFile(ctx context.Context, file models.File) error {
	log.Printf("File received, %+v", file)

	return nil
}

func NewNoneStorage() *NoneStorage {
	return &NoneStorage{}
}
