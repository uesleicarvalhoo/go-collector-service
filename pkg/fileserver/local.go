package fileserver

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type LocalFileServer struct {
	config Config
}

func NewLocalFileServer(cfg Config) (*LocalFileServer, error) {
	return &LocalFileServer{config: cfg}, nil
}

func (fs *LocalFileServer) Glob(ctx context.Context, pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func (fs *LocalFileServer) Open(ctx context.Context, filePath string) (io.ReadSeekCloser, error) {
	return os.Open(filePath)
}

func (fs *LocalFileServer) Remove(ctx context.Context, filePath string) error {
	return os.Remove(filePath)
}
