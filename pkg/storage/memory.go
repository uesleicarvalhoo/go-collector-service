package storage

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"sync"
)

var ErrFileKeyNotFound = errors.New("fileKey not found")

type MemoryStorage struct {
	storedFiles map[string][]byte
	sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		storedFiles: make(map[string][]byte),
	}
}

func (ms *MemoryStorage) SendFile(ctx context.Context, fileKey string, reader io.ReadSeeker) error {
	ms.Lock()
	defer ms.Unlock()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	ms.storedFiles[fileKey] = data

	return nil
}

func (ms *MemoryStorage) GetFile(fileKey string) ([]byte, error) {
	ms.Lock()
	defer ms.Unlock()

	if data, ok := ms.storedFiles[fileKey]; ok {
		return data, nil
	}

	return []byte{}, ErrFileKeyNotFound
}

func (ms *MemoryStorage) FileExists(fileKey string) bool {
	ms.Lock()
	defer ms.Unlock()

	_, ok := ms.storedFiles[fileKey]

	return ok
}

func (ms *MemoryStorage) RemoveFile(fileKey string) error {
	ms.Lock()
	defer ms.Unlock()

	if _, ok := ms.storedFiles[fileKey]; !ok {
		return ErrFileKeyNotFound
	}

	delete(ms.storedFiles, fileKey)

	return nil
}

func (ms *MemoryStorage) GetAllFiles() map[string][]byte {
	ms.Lock()
	defer ms.Unlock()

	return ms.storedFiles
}
