package headers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	// Adjust if your CreateOriginInfoHeader function is in another package

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert" // Adjust if your custom errors are in a different package
)

func createTestContextWithHeader(t *testing.T, headerKey string, headerValue string) *gin.Context {
	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()
	// Create a test context with the ResponseRecorder
	ctx, _ := gin.CreateTestContext(rr)

	ctx.Request = &http.Request{
		Header: http.Header{},
	}
	ctx.Request.Header.Add(headerKey, headerValue)

	// Debug: Print the context
	if ctx == nil {
		t.Fatal("Context is nil")
	}
	return ctx
}

func TestCreateOriginInfoHeaderValidHeader(t *testing.T) {
	validHeader := `{"origin":"testing","version":"1.0.0"}`
	ctx := createTestContextWithHeader(t, "X-Origin-Info", validHeader)
	originInfo, err := CreateOriginInfoHeader(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "testing", originInfo.Origin)
	assert.Equal(t, "1.0.0", originInfo.Version)
}

func TestCreateOriginInfoHeaderMissing(t *testing.T) {
	ctx := createTestContextWithHeader(t, "X-Origin-Info", "")
	_, err := CreateOriginInfoHeader(ctx)
	assert.ErrorContains(t, err, "Header schema validation")
}
