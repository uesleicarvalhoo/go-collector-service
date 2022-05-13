package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
)

func TestNewFileShoudReturnErrWhenTopicIsInvalid(t *testing.T) {
	// Arrange
	fileName := ""
	filePath := "somedir/filename.json"
	fileKey := "my-file-key"

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey)
	assert.NotNil(t, err)

	// Assert
	expectedMessage := "fileName should be informed"
	assert.Contains(t, err.Error(), expectedMessage)
}

func TestNewFileShoudReturnErrWhenKeyIsInvalid(t *testing.T) {
	// Arrange
	fileName := "filename.txt"
	filePath := ""
	fileKey := "my-file-key"

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey)
	assert.NotNil(t, err)

	// Assert
	expectedMessage := "filePath should be informed"
	assert.Contains(t, err.Error(), expectedMessage)
}

func TestNewFileShoudReturnErrWhenFileKeyIsInvalid(t *testing.T) {
	// Arrange
	fileName := "filename.txt"
	filePath := "somedir/filename.txt"
	fileKey := ""

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey)
	assert.NotNil(t, err)

	// Assert
	expectedMessage := "fileKey should be informed"
	assert.Contains(t, err.Error(), expectedMessage)
}

func TestNewFileShouldReturnAllErrorMessages(t *testing.T) {
	// Arrange
	fileName := ""
	filePath := ""
	fileKey := ""

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey)
	assert.NotNil(t, err)

	// Assert
	assert.Contains(t, err.Error(), "fileName should be informed")
	assert.Contains(t, err.Error(), "filePath should be informed")
	assert.Contains(t, err.Error(), "fileKey should be informed")
}

func TestNewFileShouldReturnFileWhenAllFieldsAreOk(t *testing.T) {
	// Arrange
	fileName := "filename.json"
	filePath := "test/filename.json"
	fileKey := "my-file-key"

	// Action
	_, err := models.NewFile(fileName, filePath, fileKey)

	// Assert
	assert.Nil(t, err)
}
