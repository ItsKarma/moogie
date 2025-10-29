package checks

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// RunDNSCheck performs a DNS resolution check
// Configuration via environment variables:
//   - DNS_HOSTNAME: Hostname to resolve (required)
//   - DNS_EXPECTED_IPS: Comma-separated list of expected IPs (optional)
//   - DNS_TIMEOUT: Timeout in seconds (default: 10)
//   - DNS_SERVER: Custom DNS server to use (optional, e.g., "8.8.8.8:53")
func RunDNSCheck() (*CheckResult, error) {
	result := NewCheckResult()

	// Parse configuration
	hostname := os.Getenv("DNS_HOSTNAME")
	if hostname == "" {
		return nil, fmt.Errorf("DNS_HOSTNAME environment variable is required")
	}

	timeout := 10
	if timeoutStr := os.Getenv("DNS_TIMEOUT"); timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = t
		}
	}

	// Create resolver
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Duration(timeout) * time.Second,
			}
			// Use custom DNS server if specified
			if dnsServer := os.Getenv("DNS_SERVER"); dnsServer != "" {
				address = dnsServer
			}
			return d.Dial(network, address)
		},
	}

	// Perform DNS lookup with timing
	ctx := context.Background()
	start := time.Now()
	ips, err := resolver.LookupHost(ctx, hostname)
	elapsed := time.Since(start)

	result.ResponseTimeMs = elapsed.Milliseconds()
	result.Metadata["hostname"] = hostname

	if err != nil {
		result.Status = "error"
		result.ErrorMessage = fmt.Sprintf("DNS lookup failed: %v", err)
		return result, nil
	}

	if len(ips) == 0 {
		result.Status = "error"
		result.ErrorMessage = "No IP addresses resolved"
		return result, nil
	}

	result.Metadata["resolved_ips"] = ips
	result.Metadata["ip_count"] = len(ips)

	// Validate expected IPs if provided
	if expectedIPsStr := os.Getenv("DNS_EXPECTED_IPS"); expectedIPsStr != "" {
		expectedIPs := strings.Split(expectedIPsStr, ",")
		expectedMap := make(map[string]bool)
		for _, ip := range expectedIPs {
			expectedMap[strings.TrimSpace(ip)] = true
		}

		// Check if at least one expected IP was resolved
		foundExpected := false
		for _, ip := range ips {
			if expectedMap[ip] {
				foundExpected = true
				break
			}
		}

		if !foundExpected {
			result.Status = "error"
			result.ErrorMessage = fmt.Sprintf("None of the expected IPs found. Expected: %v, Got: %v", expectedIPs, ips)
			return result, nil
		}
	}

	result.Status = "success"
	return result, nil
}
