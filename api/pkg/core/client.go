package core

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	StatusPending     string = "PENDING"
	StatusFail        string = "FAIL"
	StatusPass        string = "PASS"
	StatusError       string = "ERROR"
	StatusUnspecified string = "UNSPECIFIED"
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
	Status    string        `json:"status"`
	Results   []CheckResult `json:"results"`
}

func ExecuteTests(ctx context.Context, testRequest TestExecutionRequest) (*CheckResponse, error) {
	log := FromContext(ctx)
	requestID := fmt.Sprintf("aeternum-v0-%s", uuid.New().String())
	log.Debugf("Running test requests [ID %s]: %s", requestID, testRequest.BaseURL)
	for _, testRequest := range testRequest.Endpoints {
		log.Debugf("Running test request: %s", testRequest.Path)
	}
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
	failedRequests := []string{}
	for i, endpoint := range testRequest.Endpoints {
		wg.Add(1)
		go func(i int, e Endpoint) {
			defer wg.Done()
			fullURL := testRequest.BaseURL + e.Path
			resp, err := client.Get(fullURL)
			actualStatus := 0

			if err != nil {
				failedRequests = append(failedRequests, e.Path)
			}
			actualStatus = resp.StatusCode
			resp.Body.Close()
			status := "FAIL"
			if actualStatus == e.ExpectedStatus {
				status = "PASS"
			} else {
				failedRequests = append(failedRequests, e.Path)
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
	var overallStatus string
	if len(failedRequests) > 0 {
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
