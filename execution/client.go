package execution

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jgfranco17/aeternum/api/logging"
)

type Status string

const (
	StatusPending     Status = "PENDING"
	StatusFail        Status = "FAIL"
	StatusPass        Status = "PASS"
	StatusError       Status = "ERROR"
	StatusUnspecified Status = "UNSPECIFIED"
)

type Endpoint struct {
	Path           string `json:"path" binding:"required"`
	ExpectedStatus int    `json:"expected_status" binding:"required"`
}

// TestExecutionRequest represents the API health check request payload.
type TestExecutionRequest struct {
	BaseURL           string     `json:"base_url" binding:"required,url"`
	Endpoints         []Endpoint `json:"endpoints" binding:"required,dive"`
	MaxTimeoutSeconds *int       `json:"max_timeout_seconds,omitempty"`
}

// CheckResult represents the result of an individual API test.
type CheckResult struct {
	Path           string `json:"path"`
	ExpectedStatus int    `json:"expected_status"`
	ActualStatus   int    `json:"actual_status"`
	StatusCode     string `json:"status"`
}

// CheckResponse represents the full response of an API check.
type CheckResponse struct {
	RequestID string        `json:"request_id"`
	BaseURL   string        `json:"base_url"`
	Status    Status        `json:"status"`
	Results   []CheckResult `json:"results"`
}

func ExecuteTests(ctx context.Context, testRequest TestExecutionRequest) (*CheckResponse, error) {
	log := logging.FromContext(ctx)
	requestID := fmt.Sprintf("aeternum-v0-%s", uuid.New().String())
	log.Debugf("Running test requests [ID %s]: %s", requestID, testRequest.BaseURL)
	var wg sync.WaitGroup
	results := make([]CheckResult, len(testRequest.Endpoints))

	// Set timeout for API requests
	var timeout int
	if testRequest.MaxTimeoutSeconds != nil {
		timeout = *testRequest.MaxTimeoutSeconds
	} else {
		timeout = 5
	}
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	requestErrors := []string{}
	failedTests := []string{}
	for i, endpoint := range testRequest.Endpoints {
		wg.Add(1)
		go func(i int, e Endpoint) {
			defer wg.Done()
			fullURL := testRequest.BaseURL + e.Path
			resp, err := client.Get(fullURL)
			actualStatus := 0

			if err != nil {
				requestErrors = append(requestErrors, err.Error())
				return
			}
			actualStatus = resp.StatusCode
			resp.Body.Close()
			status := "FAIL"
			if actualStatus == e.ExpectedStatus {
				status = "PASS"
			} else {
				failedTests = append(failedTests, e.Path)
			}

			results[i] = CheckResult{
				Path:           e.Path,
				ExpectedStatus: e.ExpectedStatus,
				ActualStatus:   actualStatus,
				StatusCode:     status,
			}
		}(i, endpoint)
	}
	wg.Wait()

	// Handle error cases
	if len(requestErrors) > 0 {
		return nil, fmt.Errorf("Failed to make %d requests: %v", len(requestErrors), requestErrors)
	}
	var overallStatus Status
	if len(failedTests) > 0 {
		overallStatus = StatusFail
	} else {
		overallStatus = StatusPass
	}
	return &CheckResponse{
		RequestID: requestID,
		BaseURL:   testRequest.BaseURL,
		Results:   results,
		Status:    overallStatus,
	}, nil
}
