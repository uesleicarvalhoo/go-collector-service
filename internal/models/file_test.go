package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/fileserver"
)

type TestFileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

func (f TestFileInfo) Name() string {
	return f.name
}

func (f TestFileInfo) Size() int64 {
	return f.size
}

func (f TestFileInfo) ModTime() time.Time {
	return f.modTime
}

func (f TestFileInfo) IsDir() bool {
	return false
}

func newController() models.FileController {
	server, err := fileserver.NewLocalFileServer(fileserver.Config{})
	if err != nil {
		panic(err)
	}

	return server
}

func TestNewFileShoudReturnErrWhenControllerIsNill(t *testing.T) {
	// Arrange
	fileName := "myfile.txt"
	filePath := "somedir/filename.json"
	fileKey := "my-file-key"
	modTime := time.Now()
	var fileSize int64 = 120

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey, fileSize, modTime, nil)
	assert.NotNil(t, err)

	// Assert
	expectedMessage := "controller: controller is required"
	assert.Contains(t, err.Error(), expectedMessage)
}

func TestNewFileShoudReturnErrWhenTopicIsInvalid(t *testing.T) {
	// Arrange
	fileName := ""
	filePath := "somedir/filename.json"
	fileKey := "my-file-key"
	modTime := time.Now()
	var fileSize int64 = 120
	controller := newController()

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey, fileSize, modTime, controller)
	assert.NotNil(t, err)

	// Assert
	expectedMessage := "name: field is required"
	assert.Contains(t, err.Error(), expectedMessage)
}

func TestNewFileShoudReturnErrWhenKeyIsInvalid(t *testing.T) {
	// Arrange
	fileName := "filename.txt"
	filePath := ""
	fileKey := "my-file-key"
	modTime := time.Now()
	var fileSize int64 = 120
	controller := newController()

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey, fileSize, modTime, controller)
	assert.NotNil(t, err)

	// Assert
	expectedMessage := "filepath: field is required"
	assert.Contains(t, err.Error(), expectedMessage)
}

func TestNewFileShoudReturnErrWhenFileKeyIsInvalid(t *testing.T) {
	// Arrange
	fileName := "filename.txt"
	filePath := "somedir/filename.txt"
	fileKey := ""
	modTime := time.Now()
	var fileSize int64 = 120
	controller := newController()

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey, fileSize, modTime, controller)
	assert.NotNil(t, err)

	// Assert
	expectedMessage := "key: field is required"
	assert.Contains(t, err.Error(), expectedMessage)
}

func TestNewFileShouldReturnAllErrorMessages(t *testing.T) {
	// Arrange
	fileName := ""
	filePath := ""
	fileKey := ""
	modTime := time.Now()
	var fileSize int64 = 120
	controller := newController()

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey, fileSize, modTime, controller)
	assert.NotNil(t, err)

	// Assert
	assert.Contains(t, err.Error(), "name: field is required")
	assert.Contains(t, err.Error(), "filepath: field is required")
	assert.Contains(t, err.Error(), "key: field is required")
}

func TestNewFileShouldReturnFileWhenAllFieldsAreOk(t *testing.T) {
	// Arrange
	fileName := "filename.json"
	filePath := "test/filename.json"
	fileKey := "my-file-key"
	modTime := time.Now()
	var fileSize int64 = 120
	controller := newController()

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey, fileSize, modTime, controller)

	// Assert
	assert.Nil(t, err)
}
