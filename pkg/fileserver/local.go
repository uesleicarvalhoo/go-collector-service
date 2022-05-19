package fileserver

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/gofrs/flock"
)

type LocalFileServer struct {
	sync.Mutex
	config Config
}

func NewLocalFileServer(cfg Config) (*LocalFileServer, error) {
	return &LocalFileServer{
		config: cfg,
	}, nil
}

func (fs *LocalFileServer) Glob(ctx context.Context, pattern string) ([]string, error) {
	files := []string{}

	matchs, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	for _, match := range matchs {
		if f, _ := os.Stat(match); !f.IsDir() {
			absFilePath, err := filepath.Abs(match)
			if err != nil {
				return nil, err
			}

			files = append(files, absFilePath)
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

func (fs *LocalFileServer) Move(ctx context.Context, oldname, newname string) error {
	dirName, _ := filepath.Split(newname)
	if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
		return err
	}

	return os.Rename(oldname, newname)
}

func (fs *LocalFileServer) AcquireLock(ctx context.Context, filePath string) (Locker, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	locker := flock.New(filePath)

	if _, err := locker.TryLock(); err != nil {
		return nil, err
	}

	return locker, nil
}

func (fs *LocalFileServer) Stat(ctx context.Context, filepath string) (fs.FileInfo, error) {
	return os.Stat(filepath)
}
