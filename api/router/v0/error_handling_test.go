package v0

import (
	"context"
	"fmt"
	"testing"

	"github.com/jgfranco17/aeternum/api/httperror"
	"github.com/jgfranco17/aeternum/api/logging"

	"github.com/go-playground/assert/v2"
)

func TestHandeHTTPError(t *testing.T) {
	testCases := []struct {
		description    string
		response       func() errorResponse
		expectedStatus int
		expectedBody   errorBody
	}{
		{
			description: "Simple input error",
			response: func() errorResponse {
				inputErr := httperror.New(context.Background(), 500, "Some error")
				return getErrorResponse(context.Background(), inputErr)
			},
			expectedStatus: 500,
			expectedBody: errorBody{
				Message: "Some error",
			},
		},
		{
			description: "HTTP error wrapped in generic error",
			response: func() errorResponse {
				inputErr := httperror.New(context.Background(), 500, "Some error")

				err := fmt.Errorf("Outer error: %w", inputErr)
				return getErrorResponse(context.Background(), err)
			},
			expectedStatus: 500,
			expectedBody: errorBody{
				Message: "Outer error: Some error",
			},
		},
		{
			description: "HTTP error wrapping generic error",
			response: func() errorResponse {
				err := fmt.Errorf("Inner error")
				inputErr := httperror.New(context.Background(), 500, "Some error: %w", err)
				return getErrorResponse(context.Background(), inputErr)
			},
			expectedStatus: 500,
			expectedBody: errorBody{
				Message: "Some error: Inner error",
			},
		},
		{
			description: "HTTP error with requestId",
			response: func() errorResponse {
				ctx := context.WithValue(context.Background(), logging.RequestId, "4dfdcc88-2f3e-41ce-9757-4144cb3974a4")

				inputErr := httperror.New(ctx, 500, "Some error")
				return getErrorResponse(context.Background(), inputErr)
			},
			expectedStatus: 500,
			expectedBody: errorBody{
				Message:   "Some error",
				RequestID: "4dfdcc88-2f3e-41ce-9757-4144cb3974a4",
			},
		},
		{
			description: "HTTP error with service version",
			response: func() errorResponse {
				ctx := context.WithValue(context.Background(), logging.Version, "1.23.5")

				inputErr := httperror.New(ctx, 500, "Some error")
				return getErrorResponse(context.Background(), inputErr)
			},
			expectedStatus: 500,
			expectedBody: errorBody{
				Message:        "Some error",
				ServiceVersion: "1.23.5",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			response := tc.response()
			assert.Equal(t, tc.expectedStatus, response.Status)
			assert.Equal(t, tc.expectedBody, response.Body)
		})
	}
}
