package collector

import (
	"log"
	"os"
	"path/filepath"

	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
)

type LocalCollector struct {
	foldersToWatch []string
}

func NewLocalCollector(foldersToWatch ...string) *LocalCollector {
	return &LocalCollector{foldersToWatch: foldersToWatch}
}

func (lc *LocalCollector) ListFiles() (fileList []models.FileInfo, err error) {
	for _, mask := range lc.foldersToWatch {
		files, err := filepath.Glob(mask)
		if err != nil {
			log.Printf("Failed to read files in directory: %s\n", err)
		}

		for _, file := range files {
			fileList = append(
				fileList,
				models.FileInfo{
					Name:     filepath.Base(file),
					FilePath: file,
				},
			)
		}
	}

	return fileList, nil
}

func (lc *LocalCollector) GetFileData(file models.FileInfo) ([]byte, error) {
	return os.ReadFile(file.FilePath)
}

func (lc *LocalCollector) RemoveFile(file models.FileInfo) error {
	return os.Remove(file.FilePath)
}
