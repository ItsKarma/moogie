package checks

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

// RunTCPCheck performs a TCP connectivity check
// Configuration via environment variables:
//   - TCP_HOST: Target host (required)
//   - TCP_PORT: Target port (required)
//   - TCP_TIMEOUT: Timeout in seconds (default: 10)
func RunTCPCheck() (*CheckResult, error) {
	result := NewCheckResult()

	// Parse configuration
	host := os.Getenv("TCP_HOST")
	if host == "" {
		return nil, fmt.Errorf("TCP_HOST environment variable is required")
	}

	portStr := os.Getenv("TCP_PORT")
	if portStr == "" {
		return nil, fmt.Errorf("TCP_PORT environment variable is required")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid TCP_PORT: %w", err)
	}

	timeout := 10
	if timeoutStr := os.Getenv("TCP_TIMEOUT"); timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = t
		}
	}

	// Attempt TCP connection with timing
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	dialer := net.Dialer{
		Timeout: time.Duration(timeout) * time.Second,
	}

	start := time.Now()
	conn, err := dialer.Dial("tcp", addr)
	elapsed := time.Since(start)

	result.ResponseTimeMs = elapsed.Milliseconds()
	result.Metadata["host"] = host
	result.Metadata["port"] = port

	if err != nil {
		result.Status = "error"
		result.ErrorMessage = fmt.Sprintf("TCP connection failed: %v", err)
		return result, nil
	}
	defer conn.Close()

	// Get local and remote addresses
	result.Metadata["local_addr"] = conn.LocalAddr().String()
	result.Metadata["remote_addr"] = conn.RemoteAddr().String()

	result.Status = "success"
	return result, nil
}
