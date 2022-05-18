package models

import (
	"context"
	"errors"
	"io"
	"strings"
)

var ErrFileIsNotLocked = errors.New("file is not locked")

type Locker interface {
	Unlock() error
}

type File struct {
	Name       string
	FilePath   string
	Key        string
	locker     Locker
	controller FileController
}

type FileController interface {
	Open(ctx context.Context, filepath string) (io.ReadSeekCloser, error)
	Move(ctx context.Context, oldpath string, newpath string) error
	AcquireLock(ctx context.Context, filepath string) (Locker, error)
}

func NewFile(fileName, filePath string, fileKey string, controller FileController) (File, error) {
	file := File{
		Name:       fileName,
		FilePath:   filePath,
		Key:        fileKey,
		controller: controller,
	}

	if err := file.validate(); err != nil {
		return File{}, err
	}

	return file, nil
}

func (f *File) validate() error {
	validator := newValidator()

	if strings.TrimSpace(f.Name) == "" {
		validator.AddError("name", "field is required")
	}

	if strings.TrimSpace(f.FilePath) == "" {
		validator.AddError("filepath", "field is required")
	}

	if strings.TrimSpace(f.Key) == "" {
		validator.AddError("key", "field is required")
	}

	if f.controller == nil {
		validator.AddError("controller", "controller is required")
	}

	if validator.HasErrors() {
		return validator.GetError()
	}

	return nil
}

func (f *File) Open(ctx context.Context) (io.ReadSeekCloser, error) {
	return f.controller.Open(ctx, f.FilePath)
}

func (f *File) Move(ctx context.Context, newPath string) error {
	return f.controller.Move(ctx, f.FilePath, newPath)
}

func (f *File) Lock(ctx context.Context) error {
	locker, err := f.controller.AcquireLock(ctx, f.FilePath)
	if err != nil {
		return err
	}

	f.locker = locker

	return nil
}

func (f *File) Unlock(ctx context.Context) error {
	if f.locker == nil {
		return ErrFileIsNotLocked
	}

	return f.locker.Unlock()
}
