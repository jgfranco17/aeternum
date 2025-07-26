package v0

import (
	"context"
	"fmt"
	"testing"

	"github.com/jgfranco17/aeternum/api/pkg/httperror"
	"github.com/jgfranco17/aeternum/api/pkg/logging"

	"github.com/go-playground/assert/v2"
)

func TestHandeInputError(t *testing.T) {
	t.Run("Simple input error", func(t *testing.T) {
		inputErr := httperror.NewInputError(context.Background(), "Some error")
		response := getErrorResponse(context.Background(), inputErr)

		assert.Equal(t, 400, response.Status)
		assert.Equal(t, errorBody{
			Message: "Some error",
		}, response.Body)

	})

	t.Run("Input error wrapped in generic error", func(t *testing.T) {
		inputErr := httperror.NewInputError(context.Background(), "Some error")

		err := fmt.Errorf("Outer error: %w", inputErr)
		response := getErrorResponse(context.Background(), err)

		assert.Equal(t, 400, response.Status)
		assert.Equal(t, errorBody{
			Message: "Outer error: Some error",
		}, response.Body)

	})

	t.Run("Input error wrapping generic error", func(t *testing.T) {
		err := fmt.Errorf("Inner error")
		inputErr := httperror.NewInputError(context.Background(), "Some error: %w", err)

		response := getErrorResponse(context.Background(), inputErr)

		assert.Equal(t, 400, response.Status)
		assert.Equal(t, errorBody{
			Message: "Some error: Inner error",
		}, response.Body)

	})

	t.Run("Input error with requestId", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), logging.RequestId, "4dfdcc88-2f3e-41ce-9757-4144cb3974a4")

		inputErr := httperror.NewInputError(ctx, "Some error")
		response := getErrorResponse(context.Background(), inputErr)

		assert.Equal(t, 400, response.Status)
		assert.Equal(t, errorBody{
			Message:   "Some error",
			RequestID: "4dfdcc88-2f3e-41ce-9757-4144cb3974a4",
		}, response.Body)

	})

	t.Run("Input error with service version", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), logging.Version, "1.23.5")

		inputErr := httperror.NewInputError(ctx, "Some error")
		response := getErrorResponse(context.Background(), inputErr)

		assert.Equal(t, 400, response.Status)
		assert.Equal(t, errorBody{
			Message:        "Some error",
			ServiceVersion: "1.23.5",
		}, response.Body)

	})

}
