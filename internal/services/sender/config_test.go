package sender

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/collector"
)

func TestValidateShouldReturnErrorWhenEventTopicIsInvalid(t *testing.T) {
	// Arrange
	sut := Config{
		Workers: 1,
		CollectorCfg: collector.Config{
			MatchPatterns:       []string{"./files/*.json"},
			MaxCollectBatchSize: 10,
			CollectDelay:        0,
		},
	}

	// Action
	err := sut.Validate()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "eventTopic: field is required")
}

func TestValidateShouldReturnErrorWhenWorkersIsInvalid(t *testing.T) {
	// Arrange
	sut := Config{
		CollectorCfg: collector.Config{
			MatchPatterns:       []string{"./files/*.json"},
			MaxCollectBatchSize: 10,
			CollectDelay:        0,
		},
		EventTopic: "event-topic",
	}

	// Action
	err := sut.Validate()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "workers: must be higher then 0")
}

func TestValidateShouldReturnErrorWhenCollectorConfigIsInvalid(t *testing.T) {
	// Arrange
	sut := Config{
		EventTopic: "event-topic",
		Workers:    1,
	}

	// Action
	err := sut.Validate()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "collector: ")
}

func TestValidateShouldReturnAllErrors(t *testing.T) {
	// Arrange
	sut := Config{}

	// Action
	err := sut.Validate()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "collector: ")
	assert.Contains(t, err.Error(), "workers: must be higher then 0")
	assert.Contains(t, err.Error(), "eventTopic: field is required")
}

func TestValidateShouldReturnNillWhenConfigIsValid(t *testing.T) {
	// Arrange
	sut := Config{
		CollectorCfg: collector.Config{
			MatchPatterns:       []string{"./files/*.json"},
			MaxCollectBatchSize: 10,
			CollectDelay:        5,
		},
		EventTopic: "event-topic",
		Workers:    1,
	}

	// Action
	err := sut.Validate()

	// Assert
	assert.Nil(t, err)
}
