package main

import (
	"fmt"
	"log"
	"os"

	"github.com/itskarma/moogie/runner/checks"
	"github.com/itskarma/moogie/runner/client"
)

func main() {
	// Read configuration from environment variables
	checkType := os.Getenv("CHECK_TYPE")
	apiURL := os.Getenv("MOOGIE_API_URL")
	if apiURL == "" {
		apiURL = "http://moogie-api:8080"
	}
	jobName := os.Getenv("JOB_NAME")

	if checkType == "" {
		log.Fatal("CHECK_TYPE environment variable is required")
	}
	if jobName == "" {
		log.Fatal("JOB_NAME environment variable is required")
	}

	// Create API client
	apiClient := client.NewClient(apiURL)

	// Execute the check based on type
	var result *checks.CheckResult
	var err error

	switch checkType {
	case "http":
		result, err = checks.RunHTTPCheck()
	case "ssl":
		result, err = checks.RunSSLCheck()
	case "dns":
		result, err = checks.RunDNSCheck()
	case "tcp":
		result, err = checks.RunTCPCheck()
	default:
		log.Fatalf("Unknown check type: %s", checkType)
	}

	if err != nil {
		log.Printf("Check execution error: %v", err)
		// Still report the failure to the API
		result = &checks.CheckResult{
			Status:       "error",
			ErrorMessage: err.Error(),
		}
	}

	// Report result to API
	if err := apiClient.ReportExecution(jobName, result); err != nil {
		log.Fatalf("Failed to report execution result: %v", err)
	}

	fmt.Printf("Check completed successfully. Status: %s, Response Time: %dms\n",
		result.Status, result.ResponseTimeMs)
}
