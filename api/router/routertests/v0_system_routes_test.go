package routertests

import (
	"net/http"
	"testing"
)

func TestCheckNotFoundEndpoint(t *testing.T) {
	// Setup the router
	testService := NewTestServer(8800).WithSystemRoutes()

	// Define the test request (for the POST /check endpoint)
	testRequest := []ExampleHttpRequest{
		{
			Method:       "GET",
			Endpoint:     "/missing",
			ExpectedCode: http.StatusNotFound,
			Payload:      "",
		},
	}
	testService.RunRequests(t, testRequest, "")
}
