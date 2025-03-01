package core

import (
	"fmt"
	"time"

	"cli/outputs"
)

type TestResult struct {
	Target         string
	Duration       time.Duration
	Timeout        *time.Duration
	HttpStatus     int
	ExpectedStatus int
	Passed         bool
}

type TestReport struct {
	Timestamp time.Time
	Results   []TestResult
}

/*
Description: Prints a summary of the test results with colorized output.

[OUT] error: Any error occurred during the summary printing
*/
func (tr *TestReport) SummarizeResults() error {
	if len(tr.Results) < 1 {
		return fmt.Errorf("No test results to print.")
	}
	fmt.Println("------------------------------")
	for _, result := range tr.Results {
		if result.Passed {
			if result.Timeout != nil {
				if !(result.Duration < *result.Timeout) {
					outputs.PrintColoredMessage("yellow", "SLOW", "%s (%vms) exceeded threshold", result.Target, result.Duration.Milliseconds())
				}
			}
			outputs.PrintColoredMessage("green", "SUCCESS", "%s (%vms) OK", result.Target, result.Duration.Milliseconds())
		} else {
			outputs.PrintColoredMessage("red", "FAILED", "Target '%s' expected HTTP status %d but got %d", result.Target, result.ExpectedStatus, result.HttpStatus)
		}
	}
	return nil
}
