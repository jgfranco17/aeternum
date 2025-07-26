package routertests

import (
	"net/http"
	"testing"
)

func TestRunTestExecutionRequestSuccess(t *testing.T) {
	// Setup the router
	testService := NewTestServer(8800).WithSystemRoutes()

	// Define the test request (for the POST /check endpoint)
	testRequest := []ExampleHttpRequest{
		{
			Method:       "POST",
			Endpoint:     "/v0/tests/run",
			ExpectedCode: http.StatusNotFound,
			Payload: `{
				"base_url": "https://example.com/api",
				"endpoints": [
					{
						"path":            "/health",
						"expected_status": 200,
					},
					{
						"path":            "/users",
						"expected_status": 200,
					}
				]
			}`,
		},
	}
	testService.RunTestRequests(t, testRequest)
}
