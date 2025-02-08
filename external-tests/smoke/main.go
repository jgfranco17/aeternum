package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	API_TIMEOUT       int = 5
	MAX_RETRIES_COUNT int = 3
)

type SmokeTestRunner struct {
	BaseURL string
	client  *http.Client
}

func (s *SmokeTestRunner) RunSmokeTests(targetEndpoints map[string]int) error {
	failedEndpoints := make([]string, len(targetEndpoints))

	for path, expectedCode := range targetEndpoints {
		url := s.BaseURL + path
		resp, err := s.client.Get(url)
		if err != nil {
			fmt.Printf("❌ ERROR: Failed to reach %s - %v\n", url, err)
			failedEndpoints = append(failedEndpoints, path)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == 503 {
			for idx := range MAX_RETRIES_COUNT {
				fmt.Printf("[%d / %d] Service is not available, retrying in 10 seconds...\n", idx+1, MAX_RETRIES_COUNT)
				time.Sleep(10 * time.Second)
			}
		}
		if resp.StatusCode == expectedCode {
			fmt.Printf("✅ PASS: %s\n", path)
		} else {
			fmt.Printf("❌ FAIL: %s (Expected: %d, Got: %d)\n", path, expectedCode, resp.StatusCode)
			failedEndpoints = append(failedEndpoints, path)
		}
	}
	if len(failedEndpoints) > 0 {
		message := fmt.Sprintf("❌ ERROR: %d endpoints failed\n", len(failedEndpoints))
		for _, endpoint := range failedEndpoints {
			message += fmt.Sprintf("\t- %s\n", endpoint)
		}
	}
	return nil
}

func NewApiTestRunner(url string, timeoutDurationSeconds int) SmokeTestRunner {
	timeout := time.Duration(timeoutDurationSeconds) * time.Second
	return SmokeTestRunner{
		BaseURL: url,
		client:  &http.Client{Timeout: timeout},
	}
}

func main() {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		fmt.Println("❌ ERROR: API_BASE_URL is not set.")
		os.Exit(1)
	}
	fmt.Printf("Running smoke tests [%s]\n", baseURL)
	var endpoints = map[string]int{
		"/healthz":      200,
		"/metrics":      200,
		"/service-info": 200,
	}

	client := NewApiTestRunner(baseURL, API_TIMEOUT)
	err := client.RunSmokeTests(endpoints)

	if err != nil {
		fmt.Printf("Smoke tests failed: %v\n", err)
		os.Exit(1)
	}
}
