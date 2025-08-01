package execution

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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
	Method         string `json:"method" binding:"required,oneof=GET POST PUT DELETE PATCH"`
	ExpectedStatus int    `json:"expected_status" binding:"required"`
}

// TargetDefinition represents the API health check request payload.
type TargetDefinition struct {
	BaseURL           string     `json:"base_url" binding:"required"`
	Endpoints         []Endpoint `json:"endpoints" binding:"required"`
	MaxTimeoutSeconds *int       `json:"max_timeout_seconds,omitempty"`
}

// CheckResult represents the result of an individual API test.
type CheckResult struct {
	Path           string `json:"path"`
	ExpectedStatus int    `json:"expected_status"`
	ActualStatus   int    `json:"actual_status"`
	StatusCode     string `json:"status"`
}

// OutputResponse represents the full response of an API check.
type OutputResponse struct {
	RequestID string        `json:"request_id"`
	BaseURL   string        `json:"base_url"`
	Status    Status        `json:"status"`
	Results   []CheckResult `json:"results"`
}

func Run(ctx context.Context, testRequest TargetDefinition) (*OutputResponse, error) {
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
	var requestErrors []error
	failedTests := []string{}
	for i, endpoint := range testRequest.Endpoints {
		wg.Add(1)
		go func(i int, e Endpoint) {
			defer wg.Done()
			fullURL, err := url.JoinPath(testRequest.BaseURL, e.Path)
			if err != nil {
				requestErrors = append(requestErrors, err)
				return
			}
			req, err := http.NewRequest(endpoint.Method, fullURL, nil)
			if err != nil {
				requestErrors = append(requestErrors, err)
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				requestErrors = append(requestErrors, err)
				return
			}
			actualStatus := resp.StatusCode
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
		consolidatedError := errors.Join(requestErrors...)
		return nil, fmt.Errorf("Failed to make %d requests: %v", len(requestErrors), consolidatedError)
	}
	var overallStatus Status
	if len(failedTests) > 0 {
		overallStatus = StatusFail
	} else {
		overallStatus = StatusPass
	}
	return &OutputResponse{
		RequestID: requestID,
		BaseURL:   testRequest.BaseURL,
		Results:   results,
		Status:    overallStatus,
	}, nil
}
