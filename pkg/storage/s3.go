package storage

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
)

type S3Storage struct {
	bucketName string
	session    *session.Session
}

func NewS3Storage(cfg Config, region string) *S3Storage {
	return &S3Storage{
		bucketName: cfg.Bucket,
		session: session.Must(session.NewSession(&aws.Config{
			Region:   aws.String(region),
			Endpoint: aws.String(cfg.Uri),
		})),
	}
}

func (svc *S3Storage) SendFile(ctx context.Context, file models.File) error {
	_, err := s3.New(svc.session).PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(svc.bucketName),
		Key:    aws.String(file.Key),
		Body:   bytes.NewReader(file.Data),
	})

	return err
}
