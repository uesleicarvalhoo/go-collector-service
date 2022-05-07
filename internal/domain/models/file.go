package models

import (
	"io"
)

type File struct {
	Name     string
	FilePath string
	manager  FileManager
}

type FileManager interface {
	GetFileReader(filePath string) (io.ReadSeekCloser, error)
	DeleteFile(filePath string) error
}

func (f *File) GetReader() (io.ReadSeekCloser, error) {
	return f.manager.GetFileReader(f.FilePath)
}

func (f *File) Delete() error {
	return f.manager.DeleteFile(f.FilePath)
}

func NewFile(name, filePath string, manager FileManager) File {
	return File{
		Name:     name,
		FilePath: filePath,
		manager:  manager,
	}
}
