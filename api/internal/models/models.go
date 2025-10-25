package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Job represents a monitoring job configuration
type Job struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	Name      string          `json:"name" gorm:"not null;uniqueIndex"`
	Type      string          `json:"type" gorm:"not null"` // e.g., "ping", "ssl", "api-health", "dns"
	Config    json.RawMessage `json:"config" gorm:"type:jsonb;not null"`
	Enabled   bool            `json:"enabled" gorm:"default:true"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`

	// Relationships
	Executions []Execution `json:"executions,omitempty" gorm:"foreignKey:JobID"`

	// Computed fields (not stored in DB)
	SuccessRate     float64 `json:"success_rate" gorm:"-"`
	LastExecution   *time.Time `json:"last_execution" gorm:"-"`
	AvgResponseTime float64 `json:"avg_response_time" gorm:"-"`
}

// Execution represents a job execution result
type Execution struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	JobID        uint            `json:"job_id" gorm:"not null;index"`
	Status       string          `json:"status" gorm:"not null"` // "success", "failure"
	ResponseTime int64           `json:"response_time"`          // in milliseconds
	Details      json.RawMessage `json:"details" gorm:"type:jsonb"`
	Timestamp    time.Time       `json:"timestamp" gorm:"not null;index"`

	// Relationships
	Job Job `json:"job,omitempty" gorm:"foreignKey:JobID"`
}

// CreateExecutionRequest represents the request body for creating a new execution
type CreateExecutionRequest struct {
	JobName      string          `json:"job_name" binding:"required"`
	Status       string          `json:"status" binding:"required,oneof=success failure"`
	ResponseTime int64           `json:"response_time"`
	Details      json.RawMessage `json:"details"`
	Timestamp    time.Time       `json:"timestamp"`
}

// DashboardSummary represents aggregated dashboard metrics
type DashboardSummary struct {
	TotalJobs       int64                    `json:"total_jobs"`
	ActiveJobs      int64                    `json:"active_jobs"`
	OverallSuccess  float64                  `json:"overall_success_rate"`
	TotalExecutions int64                    `json:"total_executions"`
	JobSummaries    []JobSummary             `json:"job_summaries"`
	RecentActivity  []Execution              `json:"recent_activity"`
	StatusBreakdown map[string]int64         `json:"status_breakdown"`
	TypeBreakdown   map[string]int64         `json:"type_breakdown"`
}

// JobSummary represents a summary of job metrics
type JobSummary struct {
	ID              uint       `json:"id"`
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	Enabled         bool       `json:"enabled"`
	SuccessRate     float64    `json:"success_rate"`
	LastExecution   *time.Time `json:"last_execution"`
	AvgResponseTime float64    `json:"avg_response_time"`
	ExecutionCount  int64      `json:"execution_count"`
}

// WebSocketMessage represents a message sent via WebSocket
type WebSocketMessage struct {
	Type string      `json:"type"` // "execution_created", "job_updated", etc.
	Data interface{} `json:"data"`
}

// TableName overrides the table name for GORM
func (Job) TableName() string {
	return "jobs"
}

func (Execution) TableName() string {
	return "executions"
}

// BeforeCreate sets the timestamp if not provided
func (e *Execution) BeforeCreate(tx *gorm.DB) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	return nil
}
