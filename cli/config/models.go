package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// TestConfig represents the top-level structure of the smoke test suite configuration file.
type TestConfig struct {
	URL       string     `yaml:"url"`
	Endpoints []Endpoint `yaml:"endpoints"`
}

/*
Description: Runs the provided test suite and returns the test report.

[IN] filePath (string): File path to the suite configuration file

[OUT] error: Any error occurred during the test run
*/
func (tc *TestConfig) Write(filePath string) error {
	data, err := yaml.Marshal(tc)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

type Endpoint struct {
	Path           string `yaml:"path"`
	ExpectedStatus int    `yaml:"expected-status"`
	Timeout        *int   `yaml:"timeout-ms,omitempty"`
}
