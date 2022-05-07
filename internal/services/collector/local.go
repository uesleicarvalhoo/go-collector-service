package collector

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
)

type LocalCollector struct {
	patterns []string
}

func NewLocalCollector(patterns ...string) (*LocalCollector, error) {
	for _, pattern := range patterns {
		if _, err := filepath.Glob(pattern); err != nil {
			return nil, err
		}
	}

	return &LocalCollector{patterns: patterns}, nil
}

func (lc *LocalCollector) GetFiles() (fileList []models.File, err error) {
	for _, mask := range lc.patterns {
		collectedFiles, _ := filepath.Glob(mask)

		logrus.Debugf("[LocalCollector] Using '%s' find %d files", mask, len(collectedFiles))

		for _, fp := range collectedFiles {
			fileName := filepath.Base(fp)
			fileList = append(
				fileList,
				models.NewFile(fileName, fp, lc),
			)
		}
	}

	return fileList, nil
}

func (lc *LocalCollector) GetFileReader(filePath string) (io.ReadSeekCloser, error) {
	return os.Open(filePath)
}

func (lc *LocalCollector) DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

func (lc *LocalCollector) newFileModel(fp string) models.File {
	return models.NewFile(filepath.Base(fp), fp, lc)
}
