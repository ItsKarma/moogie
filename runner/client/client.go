package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/itskarma/moogie/runner/checks"
)

// Client is the API client for reporting execution results
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ExecutionRequest represents the payload to create an execution
type ExecutionRequest struct {
	JobName      string                 `json:"job_name"`
	Status       string                 `json:"status"`
	ResponseTime int64                  `json:"response_time"`
	Timestamp    time.Time              `json:"timestamp"`
	Details      map[string]interface{} `json:"details,omitempty"`
}

// ReportExecution reports a check execution result to the API
func (c *Client) ReportExecution(jobName string, result *checks.CheckResult) error {
	// Map "error" status to "failure" for API compatibility
	status := result.Status
	if status == "error" {
		status = "failure"
	}

	// Merge error message into metadata/details
	details := result.Metadata
	if details == nil {
		details = make(map[string]interface{})
	}
	if result.ErrorMessage != "" {
		details["error"] = result.ErrorMessage
	}

	payload := ExecutionRequest{
		JobName:      jobName,
		Status:       status,
		ResponseTime: result.ResponseTimeMs,
		Timestamp:    result.Timestamp,
		Details:      details,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal execution request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/executions", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API returned non-success status: %d", resp.StatusCode)
	}

	return nil
}
