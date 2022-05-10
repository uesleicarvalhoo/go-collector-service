package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
			Endpoint: aws.String(cfg.URL),
		})),
	}
}

func (svc *S3Storage) SendFile(ctx context.Context, fileKey string, reader io.ReadSeeker) error {
	_, err := s3.New(svc.session).PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(svc.bucketName),
		Key:    aws.String(fileKey),
		Body:   reader,
	})

	return err
}
