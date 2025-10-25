package services

import (
	"fmt"
	"time"

	"github.com/itskarma/moogie/api/internal/models"
	"gorm.io/gorm"
)

type ExecutionService struct {
	db         *gorm.DB
	jobService *JobService
}

func NewExecutionService(db *gorm.DB, jobService *JobService) *ExecutionService {
	return &ExecutionService{
		db:         db,
		jobService: jobService,
	}
}

// CreateExecution creates a new execution record
func (s *ExecutionService) CreateExecution(req *models.CreateExecutionRequest) (*models.Execution, error) {
	// Find the job by name
	job, err := s.jobService.GetJobByName(req.JobName)
	if err != nil {
		return nil, err
	}

	// Create the execution
	execution := &models.Execution{
		JobID:        job.ID,
		Status:       req.Status,
		ResponseTime: req.ResponseTime,
		Details:      req.Details,
		Timestamp:    req.Timestamp,
	}

	// If timestamp is not provided, use current time
	if execution.Timestamp.IsZero() {
		execution.Timestamp = time.Now()
	}

	if err := s.db.Create(execution).Error; err != nil {
		return nil, fmt.Errorf("failed to create execution: %w", err)
	}

	// Load the job relationship
	if err := s.db.Preload("Job").First(execution, execution.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load execution with job: %w", err)
	}

	return execution, nil
}

// GetExecutionsByJobID retrieves executions for a specific job
func (s *ExecutionService) GetExecutionsByJobID(jobID uint, from, to time.Time, limit int) ([]models.Execution, error) {
	var executions []models.Execution

	query := s.db.Where("job_id = ? AND timestamp BETWEEN ? AND ?", jobID, from, to).
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&executions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch executions: %w", err)
	}

	return executions, nil
}

// GetRecentExecutions retrieves the most recent executions across all jobs
func (s *ExecutionService) GetRecentExecutions(limit int) ([]models.Execution, error) {
	var executions []models.Execution

	query := s.db.Preload("Job").Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&executions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch recent executions: %w", err)
	}

	return executions, nil
}
