package collector

import (
	"io"
	"os"
	"path/filepath"

	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
)

type LocalCollector struct {
	config Config
}

// Use filepath.Glob to check all patterns, if any is invalid return error.
func NewLocalCollector(cfg Config) (*LocalCollector, error) {
	if _, err := filepath.Glob(cfg.MatchPattern); err != nil {
		return nil, err
	}

	return &LocalCollector{config: cfg}, nil
}

func (lc *LocalCollector) GetFiles() (fileList []models.File, err error) {
	collectedFiles, _ := filepath.Glob(lc.config.MatchPattern)
	logger.Debugf("[LocalCollector] Pattern '%s' find %d files\n", lc.config.MatchPattern, len(collectedFiles))

	for _, fp := range collectedFiles {
		fileName := filepath.Base(fp)
		fileList = append(
			fileList,
			models.NewFile(fileName, fp, lc),
		)
	}

	return fileList, nil
}

func (lc *LocalCollector) GetFileReader(filePath string) (io.ReadSeekCloser, error) {
	return os.Open(filePath)
}

func (lc *LocalCollector) RemoveFile(file models.File) error {
	return os.Remove(file.FilePath)
}

func (lc *LocalCollector) newFileModel(fp string) models.File {
	return models.NewFile(filepath.Base(fp), fp, lc)
}
