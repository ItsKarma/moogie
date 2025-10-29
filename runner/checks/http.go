package checks

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// RunHTTPCheck performs an HTTP/HTTPS synthetic check
// Configuration via environment variables:
//   - HTTP_URL: Target URL (required)
//   - HTTP_METHOD: HTTP method (default: GET)
//   - HTTP_TIMEOUT: Timeout in seconds (default: 30)
//   - HTTP_EXPECTED_STATUS: Expected status code (default: 200)
//   - HTTP_HEADERS: Comma-separated key:value pairs (e.g., "Authorization:Bearer token,Accept:application/json")
//   - HTTP_BODY: Request body for POST/PUT requests
func RunHTTPCheck() (*CheckResult, error) {
	result := NewCheckResult()

	// Parse configuration
	url := os.Getenv("HTTP_URL")
	if url == "" {
		return nil, fmt.Errorf("HTTP_URL environment variable is required")
	}

	method := os.Getenv("HTTP_METHOD")
	if method == "" {
		method = "GET"
	}

	timeout := 30
	if timeoutStr := os.Getenv("HTTP_TIMEOUT"); timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = t
		}
	}

	expectedStatus := 200
	if statusStr := os.Getenv("HTTP_EXPECTED_STATUS"); statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			expectedStatus = s
		}
	}

	// Build HTTP request
	var body io.Reader
	if bodyStr := os.Getenv("HTTP_BODY"); bodyStr != "" {
		body = strings.NewReader(bodyStr)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	if headersStr := os.Getenv("HTTP_HEADERS"); headersStr != "" {
		headers := strings.Split(headersStr, ",")
		for _, header := range headers {
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}

	// Execute request with timing
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	result.ResponseTimeMs = elapsed.Milliseconds()

	if err != nil {
		result.Status = "error"
		result.ErrorMessage = fmt.Sprintf("HTTP request failed: %v", err)
		return result, nil
	}
	defer resp.Body.Close()

	// Store response metadata
	result.Metadata["status_code"] = resp.StatusCode
	result.Metadata["url"] = url
	result.Metadata["method"] = method

	// Check if status code matches expectation
	if resp.StatusCode != expectedStatus {
		result.Status = "error"
		result.ErrorMessage = fmt.Sprintf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
		return result, nil
	}

	result.Status = "success"
	return result, nil
}
