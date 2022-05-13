package sender

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateConfigShouldReturnErrorWhenParalelUploadsIsLowerThen1(t *testing.T) {
	// Arrange
	sut := Config{ParalelUploads: 0}

	// Action
	err := validateConfig(sut)

	// Assert
	expectedMessage := "'ParalelUploads' must be higher then 0"
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedMessage)
}

func TestValidateConfigShouldReturnErrorWhenMatchPatternsIsEmpty(t *testing.T) {
	// Arrange
	sut := Config{
		ParalelUploads: 1,
		MatchPatterns:  []string{},
	}

	// Action
	err := validateConfig(sut)

	// Assert
	expectedMessage := "'MatchPatterns' must be have one or more patterns"
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedMessage)
}

func TestValidateConfigShouldReturnErrorAllErrors(t *testing.T) {
	// Arrange
	sut := Config{
		ParalelUploads: 0,
		MatchPatterns:  []string{},
	}

	// Action
	err := validateConfig(sut)

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "'ParalelUploads' must be higher then 0")
	assert.Contains(t, err.Error(), "'MatchPatterns' must be have one or more patterns")
}

func TestValidateConfigShouldReturnNillWhenConfigIsValid(t *testing.T) {
	// Arrange
	sut := Config{
		ParalelUploads: 1,
		MatchPatterns:  []string{"./test-files/*.json"},
	}

	// Action
	err := validateConfig(sut)

	// Assert
	assert.Nil(t, err)
}
