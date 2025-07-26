package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	exec "github.com/jgfranco17/aeternum/execution"
)

// ApiClient defines a client to interact with Aeternum API
type ApiClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewApiClient initializes and returns a new API client
func NewApiClient(baseURL string, timeoutSeconds int) *ApiClient {
	return &ApiClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

// Request sends an HTTP request to the API
func (c *ApiClient) Request(method, endpoint string, payload interface{}) error {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	// Convert payload to JSON if not nil
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to encode request body: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	// Create the HTTP request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error: %s (status: %d)", string(respData), resp.StatusCode)
	}

	// Decode JSON response if response struct is provided
	var response exec.CheckResponse
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return fmt.Errorf("Failed to decode response JSON: %w", err)
	}

	return nil
}
