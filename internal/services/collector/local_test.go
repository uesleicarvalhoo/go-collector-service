package collector

import (
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
)

var tmpDir string

func init() {
	folder, err := ioutil.TempDir("", "*")
	if err != nil {
		panic(err)
	}

	tmpDir = folder
}

func newSut(masks ...string) LocalCollector {
	patterns := []string{}
	for _, mask := range masks {
		patterns = append(patterns, filepath.Join(tmpDir, mask))
	}

	collector, err := NewLocalCollector(patterns...)
	if err != nil {
		panic(err)
	}

	return *collector
}

func createTempFile(fileName string) (models.File, error) {
	lc := newSut()
	fp := filepath.Join(tmpDir, fileName)

	err := ioutil.WriteFile(fp, []byte{}, fs.ModeAppend)
	if err != nil {
		return models.File{}, err
	}

	return lc.newFileModel(fp), nil
}

func TestGetFiles(t *testing.T) {
	// Prepare
	sut := newSut("*.json")

	// Arrange
	expectedFile, err := createTempFile("test_json_file.json")
	assert.Nil(t, err)

	ignoredFile, err := createTempFile("test_ignore_this_file.txt")
	assert.Nil(t, err)

	// Action
	collectedFiles, err := sut.GetFiles()
	assert.Nil(t, err)

	// Assert
	expectedFileIsCollected := false
	ignoredFileIsCollected := false

	for _, file := range collectedFiles {
		if file.FilePath == expectedFile.FilePath {
			expectedFileIsCollected = true
		}
		if file.FilePath == ignoredFile.FilePath {
			ignoredFileIsCollected = true
		}
	}

	assert.Truef(t, expectedFileIsCollected, "File '%s' not collected", expectedFile.Name)
	assert.Falsef(t, ignoredFileIsCollected, "File '%s' is collected", ignoredFile.Name)
}

func TestGetFilesFileDelete(t *testing.T) {
	// Prepare
	sut := newSut("*.json")

	// Arrange
	_, err := createTempFile("test_file_delete.json")
	assert.Nil(t, err)

	collectedFiles, err := sut.GetFiles()
	assert.Nil(t, err)

	collectedFile := collectedFiles[0]

	assert.FileExists(t, collectedFile.FilePath)

	// Action
	collectedFile.Delete()

	// Assert
	assert.NoFileExists(t, collectedFile.FilePath)
}

func TestFileDeleteRemoveFileWithSuccess(t *testing.T) {
	// Arrange
	file, err := createTempFile("removefilewithsuccess.json")
	assert.Nil(t, err)

	// Action
	err = file.Delete()

	// Assert
	assert.Nil(t, err)
}

func TestFileDeleteRemoveWithInvalidFilePathReturnError(t *testing.T) {
	// Prepare
	sut := newSut()

	// Arrange
	file := sut.newFileModel("./invalidpath.json")

	// Action
	err := file.Delete()

	// Assert
	assert.NotNil(t, err)
}

func TestFileGetReaderWithInvalidFilePathReturnError(t *testing.T) {
	// Arrange
	file, err := createTempFile("getreaderwithinvalidpath.json")
	assert.Nil(t, err)
	// Action
	_, err = file.GetReader()

	// Assert
	assert.NotNil(t, err)
}

func TestFileGetReaderReturnValidContentReader(t *testing.T) {
	// Prepare
	sut := newSut()
	fileData := "testing!"
	fp := filepath.Join(tmpDir, "test_file_get_reader_return_valid_content_reader.json")

	// Arrange
	file := sut.newFileModel(fp)

	err := ioutil.WriteFile(file.FilePath, []byte(fileData), 0o644)
	assert.Nil(t, err)

	// Action
	reader, err := file.GetReader()
	assert.Nil(t, err)

	// Assert
	data, err := ioutil.ReadAll(reader)
	assert.Nil(t, err)

	assert.Equal(t, []byte(fileData), data)
}
