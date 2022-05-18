package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateConfigShouldReturnErrorWhenMatchPatternIsEmpty(t *testing.T) {
	// Arrange
	sut := Config{MatchPatterns: []string{}}

	// Action
	err := sut.Validate()

	// Assert
	expectedMessage := "MatchPatterns: field is required"
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedMessage)
}

func TestValidateConfigShouldReturnErrorWhenAnyPatternIsNull(t *testing.T) {
	// Arrange
	sut := Config{MatchPatterns: []string{"", "./files/*.txt"}}

	// Action
	err := sut.Validate()

	// Assert
	expectedMessage := "MatchPattern[0]: field is required"
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedMessage)
}

func TestValidateConfigShouldReturnNillWhenConfigIsValid(t *testing.T) {
	// Arrange
	sut := Config{
		MatchPatterns: []string{"./test-files/*.json"},
	}

	// Action
	err := sut.Validate()

	// Assert
	assert.Nil(t, err)
}
