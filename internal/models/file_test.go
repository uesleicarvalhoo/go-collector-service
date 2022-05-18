package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/fileserver"
)

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

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey, nil)
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
	controller := newController()

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey, controller)
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
	controller := newController()

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey, controller)
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
	controller := newController()

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey, controller)
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
	controller := newController()

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey, controller)
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
	controller := newController()

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey, controller)

	// Assert
	assert.Nil(t, err)
}
