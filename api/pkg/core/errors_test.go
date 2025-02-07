package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputErrorNewSimpleError(t *testing.T) {
	const inputErrorMessage string = "This is an input error"

	err := NewInputError(context.Background(), inputErrorMessage)

	var expectedError InputError
	assert.ErrorAs(t, err, &expectedError)
	assert.Equal(t, inputErrorMessage, err.Error())
}

func TestInputErrorNewWrappedError(t *testing.T) {
	const rootMessage string = "This is the root"
	inputErrorMessage := "This is an input error: %v"

	rootError := fmt.Errorf(rootMessage)

	err := NewInputError(context.Background(), inputErrorMessage, rootError)

	var expectedError InputError
	assert.ErrorAs(t, err, &expectedError)
	assert.Equal(t, "This is an input error: This is the root", err.Error())
}
