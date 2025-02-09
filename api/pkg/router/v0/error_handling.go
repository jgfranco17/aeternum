package v0

import (
	"context"
	"errors"

	core_errors "api/pkg/core/errors"
	"api/pkg/core/obs"

	"github.com/gin-gonic/gin"
)

type errorBody struct {
	Message        string `json:"message,omitempty"`
	RequestID      string `json:"requestId,omitempty"`
	ServiceVersion string `json:"serviceVersion,omitempty"`
}

type errorResponse struct {
	Status int
	Body   errorBody
}

func getContextField(ctx context.Context, fieldName string) string {
	value, ok := ctx.Value(fieldName).(string)
	if !ok {
		return ""
	}
	return value
}

func getErrorMetadataFromContext(ctx context.Context) errorBody {
	requestId := getContextField(ctx, obs.RequestId)
	serviceVersion := getContextField(ctx, obs.Version)

	return errorBody{
		RequestID:      requestId,
		ServiceVersion: serviceVersion,
	}
}

// Get an error response from a core error.
// Currently the messages are generic.
func getErrorResponse(ctx context.Context, err error) errorResponse {

	errorMessage := err.Error()

	var inputErr core_errors.InputError
	if errors.As(err, &inputErr) {
		body := getErrorMetadataFromContext(inputErr.Context())
		body.Message = errorMessage
		return errorResponse{Status: 400, Body: body}
	}
	body := getErrorMetadataFromContext(ctx)
	body.Message = "Internal Server Error"
	return errorResponse{Status: 500, Body: body}

}

// Generic error handling
func handleError(c *gin.Context, err error) {
	log := obs.GetLoggerFromContext(c)
	log.Error(err)
	errorResponse := getErrorResponse(c, err)
	c.JSON(errorResponse.Status, errorResponse.Body)
}

func ServeError(c *gin.Context, status int, message string, err error) {
	c.JSON(status, gin.H{
		"status":    status,
		"message":   message,
		"traceback": err.Error(),
	})
}

// Wrapper for handlers that return errors
func WithErrorHandling(handler func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := handler(c)
		if err != nil {
			handleError(c, err)
		}
	}
}
