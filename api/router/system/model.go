package system

import (
	"encoding/json"
	"fmt"
	"os"
)

type HealthStatus struct {
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
}

type ProjectCodebase struct {
	Author       string   `json:"author"`
	Repository   string   `json:"repository"`
	Version      string   `json:"version"`
	Contributors []string `json:"contributors,omitempty"`
	License      string   `json:"license"`
	Languages    []string `json:"languages"`
	Active       bool     `json:"active"`
}

type ServiceInfo struct {
	Name        string          `json:"name"`
	Environment string          `json:"environment"`
	Uptime      string          `json:"uptime"`
	Codebase    ProjectCodebase `json:"codebase"`
}

// Reads a JSON file and unmarshals it
func GetCodebaseSpecFromFile(filePath string) (*ProjectCodebase, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()

	var data ProjectCodebase
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("Failed to parse JSON: %w", err)
	}

	return &data, nil
}

type BasicErrorInfo struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
