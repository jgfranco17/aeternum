package core

type HealthStatus struct {
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
}

type AboutInfo struct {
	Name         string   `json:"name"`
	Author       string   `json:"author"`
	Repository   string   `json:"repository"`
	Contributors []string `json:"contributors,omitempty"`
	Environment  string   `json:"environment"`
	License      string   `json:"license"`
	Languages    []string `json:"languages"`
}

type BasicErrorInfo struct {
	StatusCode int
	Message    string
}
