package fileserver

import (
	"context"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tmpDir string

func init() {
	folder, err := ioutil.TempDir("", "*")
	if err != nil {
		panic(err)
	}

	tmpDir = folder
}

func newSut() LocalFileServer {
	collector, err := NewLocalFileServer(Config{})
	if err != nil {
		panic(err)
	}

	return *collector
}

func createTempFile(fileName string) (string, error) {
	fp := filepath.Join(tmpDir, fileName)

	err := ioutil.WriteFile(fp, []byte{}, fs.ModeAppend)
	if err != nil {
		return "", err
	}

	return fp, nil
}

func TestGlob(t *testing.T) {
	// Prepare
	sut := newSut()

	// Arrange
	expectedFile, err := createTempFile("test_json_file.json")
	assert.Nil(t, err)

	ignoredFile, err := createTempFile("test_ignore_this_file.txt")
	assert.Nil(t, err)

	// Action
	pattern := filepath.Join(tmpDir, "*.json")
	collectedFiles, err := sut.Glob(context.TODO(), pattern)
	assert.Nil(t, err)

	// Assert
	expectedFileIsCollected := false
	ignoredFileIsCollected := false

	for _, file := range collectedFiles {
		if file == expectedFile {
			expectedFileIsCollected = true
		}
		if file == ignoredFile {
			ignoredFileIsCollected = true
		}
	}

	assert.Truef(t, expectedFileIsCollected, "File '%s' not collected", expectedFile)
	assert.Falsef(t, ignoredFileIsCollected, "File '%s' is collected", ignoredFile)
}

func TestRemoveFileDeleteFile(t *testing.T) {
	// Prepare
	sut := newSut()
	fileName := "test_remove_file_delete_from_storage.json"

	// Arrange
	file, err := createTempFile(fileName)
	assert.Nil(t, err)
	assert.FileExists(t, file)

	// Action
	err = sut.Remove(context.TODO(), file)
	assert.Nil(t, err)

	// Assert
	assert.NoFileExists(t, file)
}

func TestOpenReturnValidContentReader(t *testing.T) {
	// Prepare
	sut := newSut()
	fileData := "testing!"
	fileName := "test_file_get_reader_return_valid_content_reader.json"

	// Arrange
	filePath := filepath.Join(tmpDir, fileName)

	err := ioutil.WriteFile(filePath, []byte(fileData), 0o644)
	assert.Nil(t, err)

	// Action
	reader, err := sut.Open(context.TODO(), filePath)
	assert.Nil(t, err)

	// Assert
	data, err := ioutil.ReadAll(reader)
	assert.Nil(t, err)

	assert.Equal(t, []byte(fileData), data)
}
