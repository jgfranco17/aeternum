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
func PingUrl(url string, count int, timeoutSeconds int) error {
	client := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
	start := time.Now()
	log.Debugf("Checking URL %s for liveness", url)
	runningTotalDuration := 0
	successfulPings := 0
	for i := 1; i < count+1; i++ {
		resp, err := client.Head(url)
		duration := time.Since(start)
		if err != nil {
			return fmt.Errorf("Failed to reach target %s: %w", url, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			outputs.PrintColoredMessage("green", "LIVE", "Target '%s' responded in %vms (%d/%d)", url, duration.Milliseconds(), i, count)
			successfulPings += 1
		} else {
			outputs.PrintColoredMessage("red", "DOWN", "Target '%s' returned HTTP status %d (%d/%d)", url, resp.StatusCode, i, count)
		}
		runningTotalDuration += int(duration.Milliseconds())
		time.Sleep(500 * time.Millisecond)
	}
	averageRequestDuration := runningTotalDuration / count
	log.Infof("Got %d of %d pings successful, average duration of %vms", successfulPings, count, averageRequestDuration)
	return nil
}
