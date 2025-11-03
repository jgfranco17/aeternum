package routertests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jgfranco17/aeternum/api/auth"
	"github.com/jgfranco17/aeternum/api/db"
	"github.com/jgfranco17/aeternum/execution"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockDBClient struct {
	mock.Mock
}

func (m *MockDBClient) StoreTestResult(ctx context.Context, userID string, result *execution.OutputResponse) error {
	args := m.Called(ctx, userID, result)
	return args.Error(0)
}

func (m *MockDBClient) GetTestResult(ctx context.Context, userID, requestID string) (*db.TestResult, error) {
	args := m.Called(ctx, userID, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.TestResult), args.Error(1)
}

func (m *MockDBClient) GetUserTestResults(ctx context.Context, userID string, limit int) ([]db.TestResult, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.TestResult), args.Error(1)
}

func TestRunRequestSuccess_ValidSubmission(t *testing.T) {
	t.Setenv("AETERNUM_JWT_SECRET", "test-secret-key")
	token, err := auth.GenerateToken("test-user-123", "test@example.com")
	require.NoError(t, err)

	client := new(MockDBClient)
	client.On("StoreTestResult", mock.Anything, "test-user-123", mock.AnythingOfType("*execution.OutputResponse")).Return(nil)
	testService := NewTestServer(8800).WithSystemRoutes().WithV0Routes(client)

	mux := http.NewServeMux()
	for _, endpoint := range []string{"/health", "/users"} {
		mux.HandleFunc(endpoint, func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})
	}
	mockServer := httptest.NewServer(mux)

	testRequest := []ExampleHttpRequest{
		{
			Method:       "POST",
			Endpoint:     "/v0/tests/run",
			ExpectedCode: http.StatusOK,
			Payload: fmt.Sprintf(`{
				"base_url": "%s",
				"endpoints": [
					{
						"path":            "/health",
						"method":         "GET",
						"expected_status": 200
					},
					{
						"path":            "/users",
						"method":         "GET",
						"expected_status": 200
					}
				]
			}`, mockServer.URL),
		},
	}

	testService.RunRequests(t, testRequest, token)
	client.AssertExpectations(t)
}

func TestRunRequestFail_RequestNotMatch(t *testing.T) {
	t.Setenv("AETERNUM_JWT_SECRET", "test-secret-key")
	token, err := auth.GenerateToken("test-user-123", "test@example.com")
	require.NoError(t, err)

	client := new(MockDBClient)
	//expectedResponse := &execution.OutputResponse{}
	client.On("StoreTestResult", mock.Anything, "test-user-123", mock.Anything).Return(nil)
	testService := NewTestServer(8800).WithSystemRoutes().WithV0Routes(client)

	mux := http.NewServeMux()
	for _, endpoint := range []string{"/health", "/users"} {
		mux.HandleFunc(endpoint, func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})
	}
	mockServer := httptest.NewServer(mux)

	testRequest := []ExampleHttpRequest{
		{
			Method:       "POST",
			Endpoint:     "/v0/tests/run",
			ExpectedCode: http.StatusOK,
			Payload: fmt.Sprintf(`{
				"base_url": "%s",
				"endpoints": [
					{
						"path":            "/health",
						"method":         "GET",
						"expected_status": 200
					},
					{
						"path":            "/users",
						"method":         "GET",
						"expected_status": 404
					}
				]
			}`, mockServer.URL),
		},
	}

	testService.RunRequests(t, testRequest, token)
	client.AssertExpectations(t)
}

func TestRunRequestFail_NoMethodProvided(t *testing.T) {
	t.Setenv("AETERNUM_JWT_SECRET", "test-secret-key")
	token, err := auth.GenerateToken("test-user-123", "test@example.com")
	require.NoError(t, err)

	mux := http.NewServeMux()
	for _, endpoint := range []string{"/home", "/healthz"} {
		mux.HandleFunc(endpoint, func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})
	}
	mockServer := httptest.NewServer(mux)

	client := new(MockDBClient)
	testService := NewTestServer(8800).WithSystemRoutes().WithV0Routes(client)

	testRequest := []ExampleHttpRequest{
		{
			Method:       "POST",
			Endpoint:     "/v0/tests/run",
			ExpectedCode: http.StatusBadRequest,
			Payload: fmt.Sprintf(`{
				"base_url": "%s",
				"endpoints": [
					{
						"path":            "/home",
						"expected_status": 200
					},
				]
			}`, mockServer.URL),
		},
	}

	testService.RunRequests(t, testRequest, token)
	client.AssertNotCalled(t, "StoreTestResult")
}
