package models

import (
	"context"
	"errors"
	"io"
)

var ErrFileIsNotLocked = errors.New("file is not locked")

type Locker interface {
	Unlock() error
}

type FileController interface {
	Open(ctx context.Context, filepath string) (io.ReadSeekCloser, error)
	Move(ctx context.Context, oldpath string, newpath string) error
	AcquireLock(ctx context.Context, filepath string) (Locker, error)
}
