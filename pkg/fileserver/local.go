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
	files := []string{}

	matchs, err := filepath.Glob(pattern)
	if err != nil {
		return []string{}, err
	}

	for _, match := range matchs {
		if f, _ := os.Stat(match); !f.IsDir() {
			files = append(files, match)
		}
	}

	return files, nil
}

func (fs *LocalFileServer) Open(ctx context.Context, filePath string) (io.ReadSeekCloser, error) {
	return os.Open(filePath)
}

func (fs *LocalFileServer) Remove(ctx context.Context, filePath string) error {
	return os.Remove(filePath)
}

func (fs *LocalFileServer) MoveFile(ctx context.Context, oldname, newname string) error {
	dirName, _ := filepath.Split(newname)
	if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
		return err
	}

	return os.Rename(oldname, newname)
}
