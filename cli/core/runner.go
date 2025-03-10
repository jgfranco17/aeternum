package core

import (
	"fmt"
	"io"
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
func PingUrl(targetUrl string, count int, timeoutSeconds int) error {
	client := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
	log.Debugf("Checking URL %s for liveness", targetUrl)
	runningTotalDuration := 0
	successfulPings := 0
	errorsCaught := []string{}
	for i := 1; i < count+1; i++ {
		start := time.Now()
		request, err := http.NewRequest("GET", targetUrl, nil)
		resp, err := client.Do(request)
		duration := time.Since(start)
		if err != nil {
			return fmt.Errorf("Failed to reach target %s: %w", targetUrl, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			outputs.PrintColoredMessage("green", "LIVE", "Target '%s' responded in %vms (%d/%d)", targetUrl, duration.Milliseconds(), i, count)
			successfulPings += 1
		} else {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				errorsCaught = append(errorsCaught, err.Error())
			} else if string(body) != "" {
				errorsCaught = append(errorsCaught, fmt.Sprintf("Got HTTP %d status: %s", resp.StatusCode, string(body)))
			} else {
				errorsCaught = append(errorsCaught, fmt.Sprintf("Unknown HTTP %d error", resp.StatusCode))
			}
			outputs.PrintColoredMessage("red", "DOWN", "Target '%s' returned HTTP status %d (%d/%d)", targetUrl, resp.StatusCode, i, count)
		}
		runningTotalDuration += int(duration.Milliseconds())
		time.Sleep(500 * time.Millisecond)
	}
	averageRequestDuration := runningTotalDuration / count
	log.Infof("Got %d of %d pings successful, average duration of %vms", successfulPings, count, averageRequestDuration)
	if len(errorsCaught) > 0 {
		fullErrorMessage := fmt.Sprintf("Encountered %d errors during ping:", len(errorsCaught))
		for _, errorMsg := range errorsCaught {
			fullErrorMessage += fmt.Sprintf("\n\t- %s", errorMsg)
		}
		return fmt.Errorf(fullErrorMessage)
	}
	return nil
}
