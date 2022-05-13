package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasErrorsShouldReturnFalseWhenHasNoErrors(t *testing.T) {
	// Arrange
	sut := newValidator()

	// Action
	result := sut.HasErrors()

	// Assert
	assert.False(t, result)
}

func TestHasErrorsShouldReturnTrueWhenHasErrors(t *testing.T) {
	// Arrange
	sut := newValidator()
	sut.AddError(ValidationErrorProps{Context: "test", Message: "err msg"})

	// Action
	result := sut.HasErrors()

	// Assert
	assert.True(t, result)
}

func TestGetErrorShouldReturnNillWhenHasNoErrors(t *testing.T) {
	// Arrange
	sut := newValidator()

	// Action
	err := sut.GetError()

	// Assert
	assert.Nil(t, err)
}

func TestGetErrorAgroupErrorMessagesByContext(t *testing.T) {
	// Prepare
	firstError := ValidationErrorProps{Context: "test", Message: "first error"}
	secondError := ValidationErrorProps{Context: "test", Message: "second error"}
	expectedErrMessage := "test: first error, second error"

	// Arrange
	sut := newValidator()
	sut.AddError(firstError)
	sut.AddError(secondError)

	// Action
	err := sut.GetError()
	assert.NotNil(t, err)

	// Assert
	assert.Equal(t, expectedErrMessage, err.Error())
}
