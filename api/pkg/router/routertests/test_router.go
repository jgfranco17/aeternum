package routertests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/pkg/router"
	"api/pkg/router/system"
	v0 "api/pkg/router/v0"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type ExampleHttpRequest struct {
	Method         string
	Endpoint       string
	ExpectedCode   int
	Payload        interface{}
	ExpectedFields map[string]interface{}
}

func NewBasicExampleRequest(method string, endpoint string, statusCode int) ExampleHttpRequest {
	return ExampleHttpRequest{
		Method:         method,
		Endpoint:       endpoint,
		ExpectedCode:   statusCode,
		Payload:        nil,
		ExpectedFields: make(map[string]interface{}),
	}
}

func (e *ExampleHttpRequest) WithPayload(payload interface{}) ExampleHttpRequest {
	e.Payload = payload
	return *e
}

type TestServer struct {
	service *router.Service
}

/*
Create a new Test Server.

Uses chaining-builder pattern to define routes.
*/
func NewTestServer(port int) *TestServer {
	baseRouter := gin.Default()
	return &TestServer{
		service: &router.Service{
			Router: baseRouter,
			Port:   port,
		},
	}
}

func (s *TestServer) WithSystemRoutes() *TestServer {
	system.SetSystemRoutes(s.service.Router, false)
	return s
}

func (s *TestServer) WithV0Routes() *TestServer {
	v0.SetRoutes(s.service.Router)
	return s
}

func (s *TestServer) RunTestRequest(t *testing.T, sampleRequest ExampleHttpRequest) {
	t.Helper()

	// Create the request with the provided method, endpoint, and body (if any)
	var reqBody *http.Request
	if sampleRequest.Payload != nil {
		bodyBytes, _ := json.Marshal(sampleRequest.Payload)
		reqBody = httptest.NewRequest(sampleRequest.Method, sampleRequest.Endpoint, bytes.NewBuffer(bodyBytes))
	} else {
		reqBody = httptest.NewRequest(sampleRequest.Method, sampleRequest.Endpoint, nil)
	}
	recorder := httptest.NewRecorder()
	s.service.Router.ServeHTTP(recorder, reqBody)
	assert.Equal(t, sampleRequest.ExpectedCode, recorder.Code, "Expected status code to match")

	// Check if the response body matches the expected fields
	var responseBody map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoErrorf(t, err, "Failed to unmarshal JSON response body")
	for key, value := range sampleRequest.ExpectedFields {
		assert.Contains(t, responseBody, key, "Response is missing key: %s", key)
		assert.Equal(t, value, responseBody[key], "Expected value for key '%s'", key)
	}
}
