package checks

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

// RunSSLCheck performs an SSL/TLS certificate check
// Configuration via environment variables:
//   - SSL_HOST: Target host (required)
//   - SSL_PORT: Target port (default: 443)
//   - SSL_TIMEOUT: Timeout in seconds (default: 15)
//   - SSL_DAYS_WARNING: Days before expiry to warn (default: 30)
//   - SSL_CHECK_CHAIN: Validate entire certificate chain (default: true)
func RunSSLCheck() (*CheckResult, error) {
	result := NewCheckResult()

	// Parse configuration
	host := os.Getenv("SSL_HOST")
	if host == "" {
		return nil, fmt.Errorf("SSL_HOST environment variable is required")
	}

	port := 443
	if portStr := os.Getenv("SSL_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	timeout := 15
	if timeoutStr := os.Getenv("SSL_TIMEOUT"); timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = t
		}
	}

	daysWarning := 30
	if daysStr := os.Getenv("SSL_DAYS_WARNING"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			daysWarning = d
		}
	}

	// Connect to server with TLS
	addr := fmt.Sprintf("%s:%d", host, port)
	dialer := &net.Dialer{
		Timeout: time.Duration(timeout) * time.Second,
	}

	start := time.Now()
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
		ServerName: host,
	})
	elapsed := time.Since(start)

	result.ResponseTimeMs = elapsed.Milliseconds()

	if err != nil {
		result.Status = "error"
		result.ErrorMessage = fmt.Sprintf("TLS connection failed: %v", err)
		return result, nil
	}
	defer conn.Close()

	// Get certificate information
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		result.Status = "error"
		result.ErrorMessage = "No certificates found"
		return result, nil
	}

	cert := certs[0] // Leaf certificate

	// Calculate days until expiry
	now := time.Now()
	daysUntilExpiry := int(cert.NotAfter.Sub(now).Hours() / 24)

	// Store certificate metadata
	result.Metadata["host"] = host
	result.Metadata["port"] = port
	result.Metadata["issuer"] = cert.Issuer.String()
	result.Metadata["subject"] = cert.Subject.String()
	result.Metadata["not_before"] = cert.NotBefore
	result.Metadata["not_after"] = cert.NotAfter
	result.Metadata["days_until_expiry"] = daysUntilExpiry
	result.Metadata["serial_number"] = cert.SerialNumber.String()

	// Check for expiry
	if now.After(cert.NotAfter) {
		result.Status = "error"
		result.ErrorMessage = fmt.Sprintf("Certificate expired on %s", cert.NotAfter.Format(time.RFC3339))
		return result, nil
	}

	if now.Before(cert.NotBefore) {
		result.Status = "error"
		result.ErrorMessage = fmt.Sprintf("Certificate not valid until %s", cert.NotBefore.Format(time.RFC3339))
		return result, nil
	}

	// Warn if expiring soon
	if daysUntilExpiry <= daysWarning {
		result.Status = "error"
		result.ErrorMessage = fmt.Sprintf("Certificate expires in %d days (warning threshold: %d days)", daysUntilExpiry, daysWarning)
		return result, nil
	}

	result.Status = "success"
	return result, nil
}
