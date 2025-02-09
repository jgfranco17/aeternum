package system

import "time"

type HealthStatus struct {
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
}

type ServiceInfo struct {
	Name         string        `json:"name"`
	Author       string        `json:"author"`
	Repository   string        `json:"repository"`
	Contributors []string      `json:"contributors,omitempty"`
	Environment  string        `json:"environment"`
	Uptime       time.Duration `json:"uptime"`
	License      string        `json:"license"`
	Languages    []string      `json:"languages"`
}

type BasicErrorInfo struct {
	StatusCode int
	Message    string
}
