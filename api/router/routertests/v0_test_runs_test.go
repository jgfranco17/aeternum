package routertests

import (
	"context"
	"net/http"
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

func (m *MockDBClient) StoreTestResult(ctx context.Context, userID string, result *execution.CheckResponse) error {
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

func TestRunTestExecutionRequestSuccess(t *testing.T) {
	t.Setenv("AETERNUM_JWT_SECRET", "test-secret-key")
	token, err := auth.GenerateToken("test-user-123", "test@example.com")
	require.NoError(t, err)

	client := new(MockDBClient)
	client.On("StoreTestResult", mock.Anything, "test-user-123", mock.AnythingOfType("*execution.CheckResponse")).Return(nil)
	testService := NewTestServer(8800).WithSystemRoutes().WithV0Routes(client)

	testRequest := []ExampleHttpRequest{
		{
			Method:       "POST",
			Endpoint:     "/v0/tests/run",
			ExpectedCode: http.StatusOK,
			Payload: `{
				"base_url": "https://example.com/api",
				"endpoints": [
					{
						"path":            "/health",
						"expected_status": 200
					},
					{
						"path":            "/users",
						"expected_status": 200
					}
				]
			}`,
		},
	}

	testService.RunRequests(t, testRequest, token)
	client.AssertExpectations(t)
}
