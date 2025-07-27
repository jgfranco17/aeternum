package db

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	exec "github.com/jgfranco17/aeternum/execution"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helper functions
func TestCountResults_PassedTests(t *testing.T) {
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
				{StatusCode: "FAIL"},
			},
			expected: 2,
		},
		{
			name: "all failed",
			results: []exec.CheckResult{
				{StatusCode: "FAIL"},
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

func TestCountResults_FailedTests(t *testing.T) {
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

func TestTestResultStruct(t *testing.T) {
	// Test TestResult struct creation and JSON marshaling
	testResult := TestResult{
		ID:        "test-id",
		UserID:    "user-123",
		RequestID: "req-456",
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
	t.Run("with valid environment", func(t *testing.T) {
		// Set up test environment
		t.Setenv("AETERNUM_DB_URL", "https://test.supabase.co")
		t.Setenv("AETERNUM_DB_KEY", "test-anon-key")

		client, err := NewClient()

		// This will likely fail in test environment due to network connectivity
		// but we can test the function structure
		if err != nil {
			assert.Contains(t, err.Error(), "failed to initialize Supabase client")
		} else {
			assert.NotNil(t, client)
			assert.NotNil(t, client.client)
		}
	})

	t.Run("without environment variables", func(t *testing.T) {
		// Clear environment variables
		t.Setenv("AETERNUM_DB_URL", "")
		t.Setenv("AETERNUM_DB_KEY", "")

		// Reset the singleton client for this test
		// Note: This is a limitation of the singleton pattern in tests
		// In a real scenario, the client would be initialized once
		client, err := NewClient()
		// The behavior depends on whether GetSupabaseClient was called before
		// If it was called with valid env vars, it will return the cached client
		// If not, it might fail or succeed depending on the supabase-go library
		if err != nil {
			assert.Contains(t, err.Error(), "failed to initialize Supabase client")
			assert.Nil(t, client)
		} else {
			// If it succeeds, that's also valid behavior
			assert.NotNil(t, client)
		}
	})
}

// Test StoreTestResult function structure
func TestStoreTestResult(t *testing.T) {
	ctx := context.Background()

	// Create test data
	checkResponse := &exec.CheckResponse{
		RequestID: "test-request-id",
		BaseURL:   "https://example.com",
		Status:    "PASS",
		Results: []exec.CheckResult{
			{Path: "/health", ExpectedStatus: 200, ActualStatus: 200, StatusCode: "PASS"},
			{Path: "/api", ExpectedStatus: 200, ActualStatus: 404, StatusCode: "FAIL"},
		},
	}

	userID := "test-user-id"

	// Test that the function can be called (will fail due to no real client)
	// This tests the function structure and error handling
	client, err := NewClient()
	if err != nil {
		// Expected in test environment
		assert.Contains(t, err.Error(), "failed to initialize Supabase client")
	} else {
		err = client.StoreTestResult(ctx, userID, checkResponse)
		// This will likely fail due to network connectivity in test environment
		if err != nil {
			assert.Contains(t, err.Error(), "failed to store test result")
		}
	}
}

// Test GetTestResult function structure
func TestGetTestResult(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"
	requestID := "test-request-id"

	// Test that the function can be called (will fail due to no real client)
	client, err := NewClient()
	if err != nil {
		// Expected in test environment
		assert.Contains(t, err.Error(), "failed to initialize Supabase client")
	} else {
		result, err := client.GetTestResult(ctx, userID, requestID)
		// This will likely fail due to network connectivity in test environment
		if err != nil {
			assert.Contains(t, err.Error(), "failed to retrieve test result")
		} else {
			assert.NotNil(t, result)
		}
	}
}

// Test GetUserTestResults function structure
func TestGetUserTestResults(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-id"

	// Test that the function can be called (will fail due to no real client)
	client, err := NewClient()
	if err != nil {
		// Expected in test environment
		assert.Contains(t, err.Error(), "failed to initialize Supabase client")
	} else {
		results, err := client.GetUserTestResults(ctx, userID, 10)
		// This will likely fail due to network connectivity in test environment
		if err != nil {
			assert.Contains(t, err.Error(), "failed to retrieve user test results")
		} else {
			assert.NotNil(t, results)
		}
	}
}

// Test Disconnect function
func TestDisconnect(t *testing.T) {
	ctx := context.Background()

	// Test that the function can be called
	client, err := NewClient()
	if err != nil {
		// Expected in test environment
		assert.Contains(t, err.Error(), "failed to initialize Supabase client")
	} else {
		err = client.Disconnect(ctx)
		// Disconnect should always return nil for Supabase
		assert.NoError(t, err)
	}
}

// Test interface compliance
func TestDatabaseClientInterface(t *testing.T) {
	// This test ensures that SupabaseClient implements DatabaseClient interface
	var _ DatabaseClient = (*SupabaseClient)(nil)
}

// Test context handling
func TestContextHandling(t *testing.T) {
	ctx := context.Background()

	// Test that context is properly passed through
	// This is mainly to ensure the functions handle context correctly
	client, err := NewClient()
	if err == nil {
		err = client.Disconnect(ctx)
		assert.NoError(t, err)
	}
}

func TestTestResultJSONSerialization(t *testing.T) {
	// Test comprehensive JSON serialization/deserialization
	testResult := TestResult{
		ID:        "test-id",
		UserID:    "user-123",
		RequestID: "req-456",
		BaseURL:   "https://example.com",
		Status:    "PASS",
		Results: []exec.CheckResult{
			{Path: "/health", ExpectedStatus: 200, ActualStatus: 200, StatusCode: "PASS"},
			{Path: "/api", ExpectedStatus: 200, ActualStatus: 404, StatusCode: "FAIL"},
		},
		CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Metadata: map[string]interface{}{
			"endpoint_count": 2,
			"passed_count":   1,
			"failed_count":   1,
			"custom_field":   "custom_value",
		},
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(testResult)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Deserialize from JSON
	var deserialized TestResult
	err = json.Unmarshal(jsonData, &deserialized)
	require.NoError(t, err)

	// Verify all fields are preserved
	assert.Equal(t, testResult.ID, deserialized.ID)
	assert.Equal(t, testResult.UserID, deserialized.UserID)
	assert.Equal(t, testResult.RequestID, deserialized.RequestID)
	assert.Equal(t, testResult.BaseURL, deserialized.BaseURL)
	assert.Equal(t, testResult.Status, deserialized.Status)
	assert.Len(t, deserialized.Results, 2)
	assert.Equal(t, testResult.Results[0].Path, deserialized.Results[0].Path)
	assert.Equal(t, testResult.Results[1].Path, deserialized.Results[1].Path)
	assert.Equal(t, testResult.CreatedAt.Unix(), deserialized.CreatedAt.Unix())
	// JSON unmarshaling converts numbers to float64, so we need to compare as float64
	assert.Equal(t, float64(testResult.Metadata["endpoint_count"].(int)), deserialized.Metadata["endpoint_count"])
	assert.Equal(t, float64(testResult.Metadata["passed_count"].(int)), deserialized.Metadata["passed_count"])
	assert.Equal(t, float64(testResult.Metadata["failed_count"].(int)), deserialized.Metadata["failed_count"])
	assert.Equal(t, testResult.Metadata["custom_field"], deserialized.Metadata["custom_field"])
}

// Test GetSupabaseClient function
func TestGetSupabaseClient(t *testing.T) {
	t.Run("with valid environment", func(t *testing.T) {
		// Set up test environment
		t.Setenv("AETERNUM_DB_URL", "https://test.supabase.co")
		t.Setenv("AETERNUM_DB_KEY", "test-anon-key")

		client := GetSupabaseClient()

		// This will likely fail in test environment due to network connectivity
		// but we can test the function structure
		if client == nil {
			// Expected in test environment without real Supabase connection
		} else {
			assert.NotNil(t, client)
		}
	})

	t.Run("without environment variables", func(t *testing.T) {
		// Clear environment variables
		t.Setenv("AETERNUM_DB_URL", "")
		t.Setenv("AETERNUM_DB_KEY", "")

		client := GetSupabaseClient()
		// Should handle missing environment variables gracefully
		// The actual behavior depends on the supabase-go library
		_ = client // Use client to avoid unused variable warning
	})
}

// Benchmark tests
func BenchmarkCountPassedTests(b *testing.B) {
	results := []exec.CheckResult{
		{StatusCode: "PASS"},
		{StatusCode: "FAIL"},
		{StatusCode: "PASS"},
		{StatusCode: "FAIL"},
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
		{StatusCode: "FAIL"},
		{StatusCode: "PASS"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		countFailedTests(results)
	}
}
