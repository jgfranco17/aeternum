package core

import (
	"fmt"
	"net/http"
	"time"

	"cli/outputs"

	log "github.com/sirupsen/logrus"
)

/*
Description: Ping a provided URL for liveness.

[IN] url (string): Target URL to ping

[IN] timeoutSeconds (int): Timeout duration for HTTP client

[OUT] error: Any error occurred during the test run
*/
func PingUrl(url string, timeoutSeconds int) error {
	client := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
	start := time.Now()
	log.Debugf("Checking URL %s for liveness", url)
	resp, err := client.Head(url)
	duration := time.Since(start)
	if err != nil {
		return fmt.Errorf("Failed to reach target %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		outputs.PrintColoredMessage("green", "LIVE", "Target %s responded in %vms", url, duration.Milliseconds())
	} else {
		outputs.PrintColoredMessage("red", "DOWN", "Target %s returned HTTP status %d", url, resp.StatusCode)
	}
	return nil
}
