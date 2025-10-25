package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/lib/pq"
)

// Job types we'll generate
var jobTypes = []string{"http", "tcp", "dns", "ssl", "ping"}

// Possible statuses for executions
var statuses = []string{"success", "failure"}

// Sample configurations for different job types
type HTTPConfig struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   map[string]string `json:"metadata"`
	Spec       HTTPSpec          `json:"spec"`
}

type HTTPSpec struct {
	URL                  string            `json:"url"`
	Method               string            `json:"method"`
	Timeout              string            `json:"timeout"`
	ExpectedStatusCode   int               `json:"expectedStatusCode"`
	Schedule             string            `json:"schedule"`
	Retries              int               `json:"retries"`
	Headers              map[string]string `json:"headers"`
	Alerts               AlertConfig       `json:"alerts"`
}

type TCPConfig struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   map[string]string `json:"metadata"`
	Spec       TCPSpec           `json:"spec"`
}

type TCPSpec struct {
	Host     string      `json:"host"`
	Port     int         `json:"port"`
	Timeout  string      `json:"timeout"`
	Schedule string      `json:"schedule"`
	Retries  int         `json:"retries"`
	Alerts   AlertConfig `json:"alerts"`
}

type DNSConfig struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   map[string]string `json:"metadata"`
	Spec       DNSSpec           `json:"spec"`
}

type DNSSpec struct {
	Domain   string      `json:"domain"`
	Type     string      `json:"type"`
	Expected string      `json:"expected"`
	Timeout  string      `json:"timeout"`
	Schedule string      `json:"schedule"`
	Retries  int         `json:"retries"`
	Alerts   AlertConfig `json:"alerts"`
}

type SSLConfig struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   map[string]string `json:"metadata"`
	Spec       SSLSpec           `json:"spec"`
}

type SSLSpec struct {
	Host             string      `json:"host"`
	Port             int         `json:"port"`
	ExpirationDays   int         `json:"expirationDays"`
	Schedule         string      `json:"schedule"`
	Retries          int         `json:"retries"`
	Alerts           AlertConfig `json:"alerts"`
}

type PingConfig struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   map[string]string `json:"metadata"`
	Spec       PingSpec          `json:"spec"`
}

type PingSpec struct {
	Host     string      `json:"host"`
	Count    int         `json:"count"`
	Timeout  string      `json:"timeout"`
	Schedule string      `json:"schedule"`
	Retries  int         `json:"retries"`
	Alerts   AlertConfig `json:"alerts"`
}

type AlertConfig struct {
	OnFailure bool   `json:"onFailure"`
	Email     string `json:"email"`
}

func main() {
	log.Println("🌱 Starting Moogie database seeder...")

	// Database connection configuration
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "moogie")
	dbPassword := getEnv("DB_PASSWORD", "moogie123")
	dbName := getEnv("DB_NAME", "moogie")

	// Connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to database successfully")

	// Seed the database
	if err := seedDatabase(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("🎉 Database seeding completed successfully!")
}

func seedDatabase(db *sql.DB) error {
	// Clear existing data first
	log.Println("🧹 Clearing existing data...")
	if err := clearData(db); err != nil {
		return fmt.Errorf("failed to clear data: %w", err)
	}

	// Generate jobs
	log.Println("📝 Creating monitoring jobs...")
	jobIDs, err := createJobs(db, 15) // Create 15 jobs
	if err != nil {
		return fmt.Errorf("failed to create jobs: %w", err)
	}

	// Generate executions for the past 30 days
	log.Println("⚡ Creating execution history...")
	if err := createExecutions(db, jobIDs, 30); err != nil {
		return fmt.Errorf("failed to create executions: %w", err)
	}

	return nil
}

func clearData(db *sql.DB) error {
	// Delete executions first (due to foreign key constraint)
	if _, err := db.Exec("DELETE FROM executions"); err != nil {
		return fmt.Errorf("failed to clear executions: %w", err)
	}

	// Delete jobs
	if _, err := db.Exec("DELETE FROM jobs"); err != nil {
		return fmt.Errorf("failed to clear jobs: %w", err)
	}

	// Reset sequences
	if _, err := db.Exec("ALTER SEQUENCE jobs_id_seq RESTART WITH 1"); err != nil {
		return fmt.Errorf("failed to reset jobs sequence: %w", err)
	}

	if _, err := db.Exec("ALTER SEQUENCE executions_id_seq RESTART WITH 1"); err != nil {
		return fmt.Errorf("failed to reset executions sequence: %w", err)
	}

	return nil
}

func createJobs(db *sql.DB, count int) ([]int, error) {
	var jobIDs []int

	for i := 0; i < count; i++ {
		jobType := jobTypes[rand.Intn(len(jobTypes))]
		name := generateJobName(jobType, i)
		config := generateJobConfig(jobType)
		enabled := rand.Float32() < 0.9 // 90% chance of being enabled

		configJSON, err := json.Marshal(config)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal config: %w", err)
		}

		var jobID int
		err = db.QueryRow(`
			INSERT INTO jobs (name, type, config, enabled, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, name, jobType, configJSON, enabled, time.Now(), time.Now()).Scan(&jobID)
		
		if err != nil {
			return nil, fmt.Errorf("failed to insert job: %w", err)
		}

		jobIDs = append(jobIDs, jobID)
		log.Printf("  ✓ Created %s job: %s (ID: %d)", jobType, name, jobID)
	}

	return jobIDs, nil
}

func createExecutions(db *sql.DB, jobIDs []int, days int) error {
	now := time.Now()
	executionCount := 0

	for _, jobID := range jobIDs {
		// Generate random execution frequency for each job (every 5-30 minutes)
		intervalMinutes := 5 + rand.Intn(25)
		interval := time.Duration(intervalMinutes) * time.Minute

		// Calculate total executions for this job over the time period
		totalDuration := time.Duration(days) * 24 * time.Hour
		executionsForJob := int(totalDuration / interval)

		// Generate success rate for this job (70-99%)
		successRate := 0.7 + rand.Float64()*0.29

		for i := 0; i < executionsForJob; i++ {
			// Calculate timestamp (working backwards from now)
			timestamp := now.Add(-time.Duration(i) * interval)

			// Determine status based on success rate
			status := "success"
			if rand.Float64() > successRate {
				status = "failure"
			}

			// Generate realistic response time
			responseTime := generateResponseTime(status)

			// Generate details based on status
			details := generateExecutionDetails(status, responseTime)
			detailsJSON, err := json.Marshal(details)
			if err != nil {
				return fmt.Errorf("failed to marshal details: %w", err)
			}

			_, err = db.Exec(`
				INSERT INTO executions (job_id, status, response_time, details, timestamp)
				VALUES ($1, $2, $3, $4, $5)
			`, jobID, status, responseTime, detailsJSON, timestamp)

			if err != nil {
				return fmt.Errorf("failed to insert execution: %w", err)
			}

			executionCount++
		}

		log.Printf("  ✓ Created %d executions for job ID %d", executionsForJob, jobID)
	}

	log.Printf("  📊 Total executions created: %d", executionCount)
	return nil
}

func generateJobName(jobType string, index int) string {
	switch jobType {
	case "http":
		services := []string{"api", "web", "auth", "payment", "notification"}
		return fmt.Sprintf("%s-health-check-%d", services[rand.Intn(len(services))], index)
	case "tcp":
		services := []string{"database", "redis", "rabbitmq", "elasticsearch", "cache"}
		return fmt.Sprintf("%s-connection-%d", services[rand.Intn(len(services))], index)
	case "dns":
		domains := []string{"primary-dns", "backup-dns", "external-dns", "internal-dns"}
		return fmt.Sprintf("%s-resolution-%d", domains[rand.Intn(len(domains))], index)
	case "ssl":
		services := []string{"web", "api", "cdn", "gateway", "proxy"}
		return fmt.Sprintf("%s-ssl-check-%d", services[rand.Intn(len(services))], index)
	case "ping":
		targets := []string{"gateway", "server", "load-balancer", "node", "instance"}
		return fmt.Sprintf("%s-ping-%d", targets[rand.Intn(len(targets))], index)
	default:
		return fmt.Sprintf("monitoring-job-%d", index)
	}
}

func generateJobConfig(jobType string) interface{} {
	gofakeit.Seed(time.Now().UnixNano())

	metadata := map[string]string{
		"environment": []string{"production", "staging", "development"}[rand.Intn(3)],
		"team":        []string{"backend", "frontend", "devops", "platform"}[rand.Intn(4)],
		"service":     gofakeit.AppName(),
	}

	alerts := AlertConfig{
		OnFailure: true,
		Email:     gofakeit.Email(),
	}

	schedules := []string{"*/5 * * * *", "*/10 * * * *", "*/15 * * * *", "0 * * * *"}
	schedule := schedules[rand.Intn(len(schedules))]

	switch jobType {
	case "http":
		return HTTPConfig{
			APIVersion: "moogie.io/v1",
			Kind:       "HttpCheck",
			Metadata:   metadata,
			Spec: HTTPSpec{
				URL:                gofakeit.URL(),
				Method:             []string{"GET", "POST", "HEAD"}[rand.Intn(3)],
				Timeout:            fmt.Sprintf("%ds", 10+rand.Intn(20)),
				ExpectedStatusCode: []int{200, 201, 204}[rand.Intn(3)],
				Schedule:           schedule,
				Retries:            1 + rand.Intn(4),
				Headers: map[string]string{
					"User-Agent": "Moogie/1.0",
					"Accept":     "application/json",
				},
				Alerts: alerts,
			},
		}
	case "tcp":
		return TCPConfig{
			APIVersion: "moogie.io/v1",
			Kind:       "TcpCheck",
			Metadata:   metadata,
			Spec: TCPSpec{
				Host:     gofakeit.DomainName(),
				Port:     []int{80, 443, 22, 3306, 5432, 6379, 9200}[rand.Intn(7)],
				Timeout:  fmt.Sprintf("%ds", 5+rand.Intn(15)),
				Schedule: schedule,
				Retries:  1 + rand.Intn(3),
				Alerts:   alerts,
			},
		}
	case "dns":
		return DNSConfig{
			APIVersion: "moogie.io/v1",
			Kind:       "DnsCheck",
			Metadata:   metadata,
			Spec: DNSSpec{
				Domain:   gofakeit.DomainName(),
				Type:     []string{"A", "AAAA", "CNAME", "MX"}[rand.Intn(4)],
				Expected: gofakeit.IPv4Address(),
				Timeout:  fmt.Sprintf("%ds", 5+rand.Intn(10)),
				Schedule: schedule,
				Retries:  1 + rand.Intn(3),
				Alerts:   alerts,
			},
		}
	case "ssl":
		return SSLConfig{
			APIVersion: "moogie.io/v1",
			Kind:       "SslCheck",
			Metadata:   metadata,
			Spec: SSLSpec{
				Host:           gofakeit.DomainName(),
				Port:           443,
				ExpirationDays: 7 + rand.Intn(23), // 7-30 days warning
				Schedule:       schedule,
				Retries:        1 + rand.Intn(3),
				Alerts:         alerts,
			},
		}
	case "ping":
		return PingConfig{
			APIVersion: "moogie.io/v1",
			Kind:       "PingCheck",
			Metadata:   metadata,
			Spec: PingSpec{
				Host:     gofakeit.IPv4Address(),
				Count:    1 + rand.Intn(4),
				Timeout:  fmt.Sprintf("%ds", 5+rand.Intn(10)),
				Schedule: schedule,
				Retries:  1 + rand.Intn(3),
				Alerts:   alerts,
			},
		}
	default:
		return map[string]interface{}{"type": jobType}
	}
}

func generateResponseTime(status string) int64 {
	if status == "success" {
		// Success: 50-500ms typically
		return int64(50 + rand.Intn(450))
	} else {
		// Failure: either timeout (5000ms+) or quick failure (100-1000ms)
		if rand.Float32() < 0.3 {
			// Timeout scenario
			return int64(5000 + rand.Intn(5000))
		} else {
			// Quick failure
			return int64(100 + rand.Intn(900))
		}
	}
}

func generateExecutionDetails(status string, responseTime int64) map[string]interface{} {
	gofakeit.Seed(time.Now().UnixNano())

	details := map[string]interface{}{
		"response_time": responseTime,
		"timestamp":     time.Now().Format(time.RFC3339),
	}

	if status == "success" {
		details["message"] = "Check completed successfully"
		details["response_code"] = []int{200, 201, 204}[rand.Intn(3)]
		details["response_headers"] = map[string]string{
			"content-type":   "application/json",
			"content-length": fmt.Sprintf("%d", 100+rand.Intn(900)),
			"server":         []string{"nginx", "apache", "cloudflare"}[rand.Intn(3)],
		}
	} else {
		errors := []string{
			"Connection timeout",
			"HTTP 500 Internal Server Error",
			"DNS resolution failed",
			"SSL certificate expired",
			"Connection refused",
			"HTTP 404 Not Found",
			"Network unreachable",
		}
		details["error"] = errors[rand.Intn(len(errors))]
		details["response_code"] = []int{0, 404, 500, 502, 503, 504}[rand.Intn(6)]
	}

	return details
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
