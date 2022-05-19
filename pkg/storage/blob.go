package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

const (
	uploadBufferSize = 2 * 1024 * 1024
	uploadMaxBuffers = 3
)

type BlobStorage struct {
	user        string
	key         string
	container   string
	credentials *azblob.SharedKeyCredential
}

func NewBlobStorage(config Config) (*BlobStorage, error) {
	credentials, err := azblob.NewSharedKeyCredential(config.User, config.Key)
	if err != nil {
		return &BlobStorage{}, err
	}

	return &BlobStorage{
		user:        config.User,
		key:         config.Key,
		container:   config.Bucket,
		credentials: credentials,
	}, nil
}

func (svc *BlobStorage) SendFile(ctx context.Context, fileKey string, reader io.ReadSeeker) error {
	fileUrl, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", svc.user, svc.container, fileKey))
	if err != nil {
		return err
	}

	blobUrl := azblob.NewBlockBlobURL(*fileUrl, azblob.NewPipeline(svc.credentials, azblob.PipelineOptions{}))
	options := azblob.UploadStreamToBlockBlobOptions{
		BufferSize: uploadBufferSize,
		MaxBuffers: uploadMaxBuffers,
	}

	_, err = azblob.UploadStreamToBlockBlob(ctx, reader, blobUrl, options)
	if err != nil {
		return err
	}

	return nil
}
