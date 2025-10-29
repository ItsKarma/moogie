package checks

import "time"

// CheckResult represents the result of a synthetic check execution
type CheckResult struct {
	Status         string    `json:"status"`          // "success" or "error"
	ResponseTimeMs int64     `json:"response_time"`   // Response time in milliseconds
	Timestamp      time.Time `json:"timestamp"`       // When the check was executed
	ErrorMessage   string    `json:"error,omitempty"` // Error message if status is "error"
	Metadata       map[string]interface{} `json:"metadata,omitempty"` // Additional check-specific data
}

// NewCheckResult creates a new CheckResult with timestamp
func NewCheckResult() *CheckResult {
	return &CheckResult{
		Timestamp: time.Now().UTC(),
		Metadata:  make(map[string]interface{}),
	}
}
