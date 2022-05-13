package models

import (
	"strings"
)

type File struct {
	Name     string
	FilePath string
	Key      string
}

func NewFile(fileName, filePath string, fileKey string) (File, error) {
	file := File{
		Name:     fileName,
		FilePath: filePath,
		Key:      fileKey,
	}

	if err := file.validate(); err != nil {
		return File{}, err
	}

	return file, nil
}

func (f *File) validate() error {
	validator := newValidator()

	if strings.TrimSpace(f.Name) == "" {
		validator.AddError(ValidationErrorProps{Context: "file", Message: "fileName should be informed"})
	}

	if strings.TrimSpace(f.FilePath) == "" {
		validator.AddError(ValidationErrorProps{Context: "file", Message: "filePath should be informed"})
	}

	if strings.TrimSpace(f.Key) == "" {
		validator.AddError(ValidationErrorProps{Context: "file", Message: "fileKey should be informed"})
	}

	if validator.HasErrors() {
		return validator.GetError()
	}

	return nil
}