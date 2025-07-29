package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jgfranco17/aeternum/api/logging"
	"github.com/jgfranco17/aeternum/execution"
	exec "github.com/jgfranco17/aeternum/execution"
	supabase "github.com/supabase-community/supabase-go"
)

// TestResult represents a stored test execution result
type TestResult struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	RequestID string                 `json:"request_id"`
	BaseURL   string                 `json:"base_url"`
	Status    execution.Status       `json:"status"`
	Results   []exec.CheckResult     `json:"results"`
	CreatedAt time.Time              `json:"created_at"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// DatabaseClient interface for database operations
type DatabaseClient interface {
	StoreTestResult(ctx context.Context, userID string, result *exec.CheckResponse) error
	GetTestResult(ctx context.Context, userID, requestID string) (*TestResult, error)
	GetUserTestResults(ctx context.Context, userID string, limit int) ([]TestResult, error)
}

// SupabaseClient implements DatabaseClient for Supabase
type SupabaseClient struct {
	client *supabase.Client
}

// NewClient creates a new Supabase database client
func NewClient() (*SupabaseClient, error) {
	client := GetSupabaseClient()
	if client == nil {
		return nil, fmt.Errorf("failed to initialize Supabase client")
	}
	return &SupabaseClient{client: client}, nil
}

// StoreTestResult stores a test execution result in Supabase
func (s *SupabaseClient) StoreTestResult(ctx context.Context, userID string, result *exec.CheckResponse) error {
	log := logging.FromContext(ctx)

	// Create the test result data
	testResult := TestResult{
		ID:        result.RequestID,
		UserID:    userID,
		RequestID: result.RequestID,
		BaseURL:   result.BaseURL,
		Status:    result.Status,
		Results:   result.Results,
		CreatedAt: time.Now(),
		Metadata: map[string]interface{}{
			"endpoint_count": len(result.Results),
			"passed_count":   countPassedTests(result.Results),
			"failed_count":   countFailedTests(result.Results),
		},
	}

	// Use Supabase SDK's ORM-like interface to insert the test result
	// Based on the documentation: client.From("table").Insert(data).Execute()
	_, count, err := s.client.From("test_results").Insert(testResult, false, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("failed to store test result: %w", err)
	}

	log.Infof("Successfully stored test result with ID: %s (count: %d)", result.RequestID, count)
	return nil
}

// GetTestResult retrieves a specific test result by request ID
func (s *SupabaseClient) GetTestResult(ctx context.Context, userID, requestID string) (*TestResult, error) {
	log := logging.FromContext(ctx)

	var results []TestResult
	data, _, err := s.client.From("test_results").
		Select("*", "exact", false).
		Eq("id", requestID).
		Eq("user_id", userID).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve test result: %w", err)
	}

	// Parse the response data
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal test result: %w", err)
	}

	if len(results) == 0 {
		log.Infof("No test result found for ID: %s", requestID)
		return nil, fmt.Errorf("test result not found")
	}

	log.Infof("Successfully retrieved test result with ID: %s", requestID)
	return &results[0], nil
}

// GetUserTestResults retrieves all test results for a user
func (s *SupabaseClient) GetUserTestResults(ctx context.Context, userID string, limit int) ([]TestResult, error) {
	log := logging.FromContext(ctx)

	// Build the query
	query := s.client.From("test_results").
		Select("*", "exact", false).
		Eq("user_id", userID)

	if limit > 0 {
		query = query.Limit(limit, "")
	}

	data, _, err := query.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user test results: %w", err)
	}

	// Parse the response data
	var results []TestResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user test results: %w", err)
	}

	log.Infof("Successfully retrieved %d test results for user: %s", len(results), userID)
	return results, nil
}

// Helper functions
func countPassedTests(results []exec.CheckResult) int {
	count := 0
	for _, result := range results {
		if result.StatusCode == "PASS" {
			count++
		}
	}
	return count
}

func countFailedTests(results []exec.CheckResult) int {
	count := 0
	for _, result := range results {
		if result.StatusCode == "FAIL" {
			count++
		}
	}
	return count
}
