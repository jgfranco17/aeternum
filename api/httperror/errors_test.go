package httperror

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpErrorNewSimpleError(t *testing.T) {
	const inputErrorMessage string = "This is an input error"

	err := New(context.Background(), 400, inputErrorMessage)

	var expectedError HttpError
	assert.ErrorAs(t, err, &expectedError)
	assert.Equal(t, inputErrorMessage, err.Error())
}

func TestHttpErrorNewWrappedError(t *testing.T) {
	const rootMessage string = "This is the root"
	inputErrorMessage := "This is an input error: %v"

	rootError := fmt.Errorf(rootMessage)

	err := New(context.Background(), 500, inputErrorMessage, rootError)

	var expectedError HttpError
	assert.ErrorAs(t, err, &expectedError)
	assert.Equal(t, "This is an input error: This is the root", err.Error())
}
