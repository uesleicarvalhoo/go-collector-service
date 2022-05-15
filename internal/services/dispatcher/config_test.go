package dispatcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/collector"
	"github.com/uesleicarvalhoo/go-collector-service/internal/services/sender"
)

func TestValidateShouldReturnErrorWhenSenderConfigIsEmpty(t *testing.T) {
	// Arrange
	sut := Config{}

	// Action
	err := sut.Validate()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "sender config is required")
}

func TestValidateShouldReturnErrorWhenSenderConfigIsInvalid(t *testing.T) {
	// Arrange
	sut := Config{
		[]sender.Config{{}},
	}

	// Action
	err := sut.Validate()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Worker[1]: ")
}

func TestValidateShouldReturnNillWhenConfigIsValid(t *testing.T) {
	// Arrange
	sut := Config{
		[]sender.Config{
			{
				EventTopic: "event-topic",
				Workers:    1,
				CollectorCfg: collector.Config{
					MaxCollectBatchSize: 10,
					CollectDelay:        5,
					MatchPatterns:       []string{"./files.json"},
				},
			},
		},
	}

	// Action
	err := sut.Validate()

	// Assert
	assert.Nil(t, err)
}
