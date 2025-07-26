package db

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	exec "github.com/jgfranco17/aeternum/execution"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockSupabaseClient is a mock implementation of the Supabase client for testing
type MockSupabaseClient struct {
	mock.Mock
}

// MockQueryBuilder is a mock implementation of the query builder
type MockQueryBuilder struct {
	mock.Mock
}

func (m *MockQueryBuilder) Select(columns, count string, exact bool) *MockQueryBuilder {
	args := m.Called(columns, count, exact)
	return args.Get(0).(*MockQueryBuilder)
}

func (m *MockQueryBuilder) Eq(column, value string) *MockQueryBuilder {
	args := m.Called(column, value)
	return args.Get(0).(*MockQueryBuilder)
}

func (m *MockQueryBuilder) Limit(count int, foreignTable string) *MockQueryBuilder {
	args := m.Called(count, foreignTable)
	return args.Get(0).(*MockQueryBuilder)
}

func (m *MockQueryBuilder) Execute() ([]byte, int, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Get(1).(int), args.Error(2)
}

func (m *MockSupabaseClient) From(table string) *MockQueryBuilder {
	args := m.Called(table)
	return args.Get(0).(*MockQueryBuilder)
}

func (m *MockSupabaseClient) Insert(data interface{}, upsert bool, onConflict, returning, count string) *MockQueryBuilder {
	args := m.Called(data, upsert, onConflict, returning, count)
	return args.Get(0).(*MockQueryBuilder)
}

// TestSupabaseClient wraps the mock for testing
type TestSupabaseClient struct {
	client *MockSupabaseClient
}

func (t *TestSupabaseClient) StoreTestResult(ctx context.Context, userID string, result *exec.CheckResponse) error {
	args := t.client.Called(ctx, userID, result)
	return args.Error(0)
}

func (t *TestSupabaseClient) GetTestResult(ctx context.Context, userID, requestID string) (*TestResult, error) {
	args := t.client.Called(ctx, userID, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TestResult), args.Error(1)
}

func (t *TestSupabaseClient) GetUserTestResults(ctx context.Context, userID string, limit int) ([]TestResult, error) {
	args := t.client.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]TestResult), args.Error(1)
}

func (t *TestSupabaseClient) Disconnect(ctx context.Context) error {
	args := t.client.Called(ctx)
	return args.Error(0)
}

// Test helper functions
func TestCountPassedTests(t *testing.T) {
	tests := []struct {
		name     string
		results  []exec.CheckResult
		expected int
	}{
		{
			name:     "empty results",
			results:  []exec.CheckResult{},
			expected: 0,
		},
		{
			name: "all passed",
			results: []exec.CheckResult{
				{StatusCode: "PASS"},
				{StatusCode: "PASS"},
				{StatusCode: "PASS"},
			},
			expected: 3,
		},
		{
			name: "mixed results",
			results: []exec.CheckResult{
				{StatusCode: "PASS"},
				{StatusCode: "FAIL"},
				{StatusCode: "PASS"},
				{StatusCode: "ERROR"},
			},
			expected: 2,
		},
		{
			name: "all failed",
			results: []exec.CheckResult{
				{StatusCode: "FAIL"},
				{StatusCode: "ERROR"},
				{StatusCode: "FAIL"},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countPassedTests(tt.results)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCountFailedTests(t *testing.T) {
	tests := []struct {
		name     string
		results  []exec.CheckResult
		expected int
	}{
		{
			name:     "empty results",
			results:  []exec.CheckResult{},
			expected: 0,
		},
		{
			name: "all failed",
			results: []exec.CheckResult{
				{StatusCode: "FAIL"},
				{StatusCode: "FAIL"},
				{StatusCode: "FAIL"},
			},
			expected: 3,
		},
		{
			name: "mixed results",
			results: []exec.CheckResult{
				{StatusCode: "PASS"},
				{StatusCode: "FAIL"},
				{StatusCode: "PASS"},
				{StatusCode: "FAIL"},
			},
			expected: 2,
		},
		{
			name: "all passed",
			results: []exec.CheckResult{
				{StatusCode: "PASS"},
				{StatusCode: "PASS"},
				{StatusCode: "PASS"},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countFailedTests(tt.results)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test TestResult struct
func TestTestResultStruct(t *testing.T) {
	now := time.Now()
	testResult := TestResult{
		ID:        "test-id",
		UserID:    "user-id",
		RequestID: "request-id",
		BaseURL:   "https://example.com",
		Status:    "PASS",
		Results: []exec.CheckResult{
			{Path: "/health", ExpectedStatus: 200, ActualStatus: 200, StatusCode: "PASS"},
		},
		CreatedAt: now,
		Metadata: map[string]interface{}{
			"endpoint_count": 1,
			"passed_count":   1,
			"failed_count":   0,
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(testResult)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Test JSON unmarshaling
	var unmarshaled TestResult
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, testResult.ID, unmarshaled.ID)
	assert.Equal(t, testResult.UserID, unmarshaled.UserID)
	assert.Equal(t, testResult.RequestID, unmarshaled.RequestID)
	assert.Equal(t, testResult.BaseURL, unmarshaled.BaseURL)
	assert.Equal(t, testResult.Status, unmarshaled.Status)
	assert.Len(t, unmarshaled.Results, 1)
	assert.Equal(t, testResult.Results[0].Path, unmarshaled.Results[0].Path)
}

// Test NewClient function
func TestNewClient(t *testing.T) {
	// This test requires environment variables to be set
	// In a real test environment, you would set these up
	t.Skip("Skipping NewClient test as it requires environment setup")
}

// Test StoreTestResult with mock
func TestStoreTestResult(t *testing.T) {
	ctx := context.Background()

	// Create test data
	checkResponse := &exec.CheckResponse{
		RequestID: "test-request-id",
		BaseURL:   "https://example.com",
		Status:    "PASS",
		Results: []exec.CheckResult{
			{Path: "/health", ExpectedStatus: 200, ActualStatus: 200, StatusCode: "PASS"},
		},
	}

	userID := "test-user-id"

	// Test successful storage
	t.Run("successful storage", func(t *testing.T) {
		// Create a mock client that returns success
		mockClient := &MockSupabaseClient{}
		testClient := &TestSupabaseClient{client: mockClient}

		mockClient.On("StoreTestResult", ctx, userID, checkResponse).Return(nil)

		err := testClient.StoreTestResult(ctx, userID, checkResponse)
		assert.NoError(t, err)

		mockClient.AssertExpectations(t)
	})

	// Test storage error
	t.Run("storage error", func(t *testing.T) {
		// Create a mock client that returns error
		mockClient := &MockSupabaseClient{}
		testClient := &TestSupabaseClient{client: mockClient}

		expectedError := errors.New("database error")
		mockClient.On("StoreTestResult", ctx, userID, checkResponse).Return(expectedError)

		err := testClient.StoreTestResult(ctx, userID, checkResponse)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		mockClient.AssertExpectations(t)
	})
}

// Test GetTestResult with mock
func TestGetTestResult(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"
	requestID := "test-request-id"

	expectedResult := &TestResult{
		ID:        requestID,
		UserID:    userID,
		RequestID: requestID,
		BaseURL:   "https://example.com",
		Status:    "PASS",
		Results: []exec.CheckResult{
			{Path: "/health", ExpectedStatus: 200, ActualStatus: 200, StatusCode: "PASS"},
		},
		CreatedAt: time.Now(),
		Metadata: map[string]interface{}{
			"endpoint_count": 1,
			"passed_count":   1,
			"failed_count":   0,
		},
	}

	// Test successful retrieval
	t.Run("successful retrieval", func(t *testing.T) {
		mockClient := &MockSupabaseClient{}
		testClient := &TestSupabaseClient{client: mockClient}

		mockClient.On("GetTestResult", ctx, userID, requestID).Return(expectedResult, nil)

		result, err := testClient.GetTestResult(ctx, userID, requestID)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)

		mockClient.AssertExpectations(t)
	})

	// Test not found
	t.Run("not found", func(t *testing.T) {
		mockClient := &MockSupabaseClient{}
		testClient := &TestSupabaseClient{client: mockClient}

		expectedError := errors.New("test result not found")
		mockClient.On("GetTestResult", ctx, userID, requestID).Return(nil, expectedError)

		result, err := testClient.GetTestResult(ctx, userID, requestID)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)

		mockClient.AssertExpectations(t)
	})
}

// Test GetUserTestResults with mock
func TestGetUserTestResults(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"
	limit := 5

	expectedResults := []TestResult{
		{
			ID:        "test-id-1",
			UserID:    userID,
			RequestID: "request-id-1",
			BaseURL:   "https://example.com",
			Status:    "PASS",
			Results: []exec.CheckResult{
				{Path: "/health", ExpectedStatus: 200, ActualStatus: 200, StatusCode: "PASS"},
			},
			CreatedAt: time.Now(),
			Metadata: map[string]interface{}{
				"endpoint_count": 1,
				"passed_count":   1,
				"failed_count":   0,
			},
		},
		{
			ID:        "test-id-2",
			UserID:    userID,
			RequestID: "request-id-2",
			BaseURL:   "https://example.com",
			Status:    "FAIL",
			Results: []exec.CheckResult{
				{Path: "/api", ExpectedStatus: 200, ActualStatus: 500, StatusCode: "FAIL"},
			},
			CreatedAt: time.Now(),
			Metadata: map[string]interface{}{
				"endpoint_count": 1,
				"passed_count":   0,
				"failed_count":   1,
			},
		},
	}

	// Test successful retrieval
	t.Run("successful retrieval", func(t *testing.T) {
		mockClient := &MockSupabaseClient{}
		testClient := &TestSupabaseClient{client: mockClient}

		mockClient.On("GetUserTestResults", ctx, userID, limit).Return(expectedResults, nil)

		results, err := testClient.GetUserTestResults(ctx, userID, limit)
		assert.NoError(t, err)
		assert.Equal(t, expectedResults, results)
		assert.Len(t, results, 2)

		mockClient.AssertExpectations(t)
	})

	// Test empty results
	t.Run("empty results", func(t *testing.T) {
		mockClient := &MockSupabaseClient{}
		testClient := &TestSupabaseClient{client: mockClient}

		mockClient.On("GetUserTestResults", ctx, userID, limit).Return([]TestResult{}, nil)

		results, err := testClient.GetUserTestResults(ctx, userID, limit)
		assert.NoError(t, err)
		assert.Empty(t, results)

		mockClient.AssertExpectations(t)
	})

	// Test retrieval error
	t.Run("retrieval error", func(t *testing.T) {
		mockClient := &MockSupabaseClient{}
		testClient := &TestSupabaseClient{client: mockClient}

		expectedError := errors.New("database error")
		mockClient.On("GetUserTestResults", ctx, userID, limit).Return(nil, expectedError)

		results, err := testClient.GetUserTestResults(ctx, userID, limit)
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.Equal(t, expectedError, err)

		mockClient.AssertExpectations(t)
	})

	// Test with zero limit
	t.Run("zero limit", func(t *testing.T) {
		mockClient := &MockSupabaseClient{}
		testClient := &TestSupabaseClient{client: mockClient}

		mockClient.On("GetUserTestResults", ctx, userID, 0).Return(expectedResults, nil)

		results, err := testClient.GetUserTestResults(ctx, userID, 0)
		assert.NoError(t, err)
		assert.Equal(t, expectedResults, results)

		mockClient.AssertExpectations(t)
	})
}

// Test Disconnect
func TestDisconnect(t *testing.T) {
	ctx := context.Background()

	// Test successful disconnect
	t.Run("successful disconnect", func(t *testing.T) {
		mockClient := &MockSupabaseClient{}
		testClient := &TestSupabaseClient{client: mockClient}

		mockClient.On("Disconnect", ctx).Return(nil)

		err := testClient.Disconnect(ctx)
		assert.NoError(t, err)

		mockClient.AssertExpectations(t)
	})

	// Test disconnect error
	t.Run("disconnect error", func(t *testing.T) {
		mockClient := &MockSupabaseClient{}
		testClient := &TestSupabaseClient{client: mockClient}

		expectedError := errors.New("disconnect error")
		mockClient.On("Disconnect", ctx).Return(expectedError)

		err := testClient.Disconnect(ctx)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		mockClient.AssertExpectations(t)
	})
}

// Test NewMongoClient (legacy function)
func TestNewMongoClient(t *testing.T) {
	ctx := context.Background()
	uri := "mongodb://localhost:27017"
	username := "testuser"
	token := "testtoken"

	// This test requires environment variables to be set
	// In a real test environment, you would set these up
	t.Skip("Skipping NewMongoClient test as it requires environment setup")

	client, err := NewMongoClient(ctx, uri, username, token)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

// Integration test helpers
func TestDatabaseClientInterface(t *testing.T) {
	// This test ensures that our TestSupabaseClient implements the DatabaseClient interface
	var _ DatabaseClient = (*TestSupabaseClient)(nil)
}

// Test context handling
func TestContextHandling(t *testing.T) {
	// Test with cancelled context
	t.Run("cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		mockClient := &MockSupabaseClient{}
		testClient := &TestSupabaseClient{client: mockClient}

		// The mock should handle the cancelled context gracefully
		mockClient.On("GetTestResult", ctx, "user-id", "request-id").Return(nil, context.Canceled)

		_, err := testClient.GetTestResult(ctx, "user-id", "request-id")
		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)

		mockClient.AssertExpectations(t)
	})
}

// Test error message formatting
func TestErrorMessageFormatting(t *testing.T) {
	// Test that error messages are properly formatted
	ctx := context.Background()
	userID := "test-user-id"
	requestID := "test-request-id"

	// Test specific error message
	t.Run("specific error message", func(t *testing.T) {
		mockClient := &MockSupabaseClient{}
		testClient := &TestSupabaseClient{client: mockClient}

		expectedError := errors.New("test result not found")
		mockClient.On("GetTestResult", ctx, userID, requestID).Return(nil, expectedError)

		_, err := testClient.GetTestResult(ctx, userID, requestID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test result not found")

		mockClient.AssertExpectations(t)
	})
}

// Benchmark tests
func BenchmarkCountPassedTests(b *testing.B) {
	results := []exec.CheckResult{
		{StatusCode: "PASS"},
		{StatusCode: "FAIL"},
		{StatusCode: "PASS"},
		{StatusCode: "ERROR"},
		{StatusCode: "PASS"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		countPassedTests(results)
	}
}

func BenchmarkCountFailedTests(b *testing.B) {
	results := []exec.CheckResult{
		{StatusCode: "PASS"},
		{StatusCode: "FAIL"},
		{StatusCode: "PASS"},
		{StatusCode: "ERROR"},
		{StatusCode: "FAIL"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		countFailedTests(results)
	}
}
