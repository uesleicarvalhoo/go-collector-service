package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-collector-service/internal/models"
)

func TestNewEventShoudReturnErrWhenTopicIsInvalid(t *testing.T) {
	// Arrange
	topic := ""
	key := "event-key"
	data := "event-data"

	// Action
	_, err := models.NewEvent(topic, key, data)
	assert.NotNil(t, err)

	// Assert
	expectedMessage := "topic should be informed"
	assert.Contains(t, err.Error(), expectedMessage)
}

func TestNewEventShoudReturnErrWhenKeyIsInvalid(t *testing.T) {
	// Arrange
	topic := "event-topic"
	key := ""
	data := "event-data"

	// Action
	_, err := models.NewEvent(topic, key, data)
	assert.NotNil(t, err)

	// Assert
	expectedMessage := "key should be informed"
	assert.Contains(t, err.Error(), expectedMessage)
}

func TestNewEventShouldReturnAllErrorMessages(t *testing.T) {
	// Arrange
	topic := "event-topic"
	key := ""
	data := "event-data"

	// Action
	_, err := models.NewEvent(topic, key, data)
	assert.NotNil(t, err)

	// Assert
	expectedMessage := "key should be informed"
	assert.Contains(t, err.Error(), expectedMessage)
}

func TestNewEventErrorShouldReturnNillWhenAllDataAreOk(t *testing.T) {
	// Arrange
	topic := "event-topic"
	key := "event-key"
	data := "event-data"

	// Action
	_, err := models.NewEvent(topic, key, data)

	// Assert
	assert.Nil(t, err)
}
