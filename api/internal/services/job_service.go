package services

import (
	"fmt"
	"time"

	"github.com/itskarma/moogie/api/internal/models"
	"gorm.io/gorm"
)

type JobService struct {
	db *gorm.DB
}

func NewJobService(db *gorm.DB) *JobService {
	return &JobService{db: db}
}

// GetAllJobs retrieves all jobs with computed metrics for the given date range
func (s *JobService) GetAllJobs(from, to time.Time) ([]models.Job, error) {
	var jobs []models.Job

	// Get all jobs
	if err := s.db.Find(&jobs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch jobs: %w", err)
	}

	// Compute metrics for each job
	for i := range jobs {
		if err := s.computeJobMetrics(&jobs[i], from, to); err != nil {
			return nil, fmt.Errorf("failed to compute metrics for job %d: %w", jobs[i].ID, err)
		}
	}

	return jobs, nil
}

// GetJobByID retrieves a job by ID with execution history and computed metrics
func (s *JobService) GetJobByID(id uint, from, to time.Time, limit int) (*models.Job, error) {
	var job models.Job

	// Get the job
	if err := s.db.First(&job, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to fetch job: %w", err)
	}

	// Get execution history for the date range
	var executions []models.Execution
	query := s.db.Where("job_id = ? AND timestamp BETWEEN ? AND ?", id, from, to).
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&executions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch executions: %w", err)
	}

	job.Executions = executions

	// Compute metrics
	if err := s.computeJobMetrics(&job, from, to); err != nil {
		return nil, fmt.Errorf("failed to compute metrics: %w", err)
	}

	return &job, nil
}

// GetJobByName retrieves a job by name
func (s *JobService) GetJobByName(name string) (*models.Job, error) {
	var job models.Job

	if err := s.db.Where("name = ?", name).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to fetch job: %w", err)
	}

	return &job, nil
}

// computeJobMetrics calculates success rate, last execution, and avg response time
func (s *JobService) computeJobMetrics(job *models.Job, from, to time.Time) error {
	// Count total executions in date range
	var totalCount int64
	if err := s.db.Model(&models.Execution{}).
		Where("job_id = ? AND timestamp BETWEEN ? AND ?", job.ID, from, to).
		Count(&totalCount).Error; err != nil {
		return err
	}

	if totalCount == 0 {
		job.SuccessRate = 0
		job.AvgResponseTime = 0
		job.LastExecution = nil
		return nil
	}

	// Count successful executions
	var successCount int64
	if err := s.db.Model(&models.Execution{}).
		Where("job_id = ? AND status = 'success' AND timestamp BETWEEN ? AND ?", job.ID, from, to).
		Count(&successCount).Error; err != nil {
		return err
	}

	// Calculate success rate
	job.SuccessRate = float64(successCount) / float64(totalCount) * 100

	// Get average response time
	var avgResponseTime float64
	if err := s.db.Model(&models.Execution{}).
		Where("job_id = ? AND timestamp BETWEEN ? AND ?", job.ID, from, to).
		Select("AVG(response_time)").
		Scan(&avgResponseTime).Error; err != nil {
		return err
	}
	job.AvgResponseTime = avgResponseTime

	// Get last execution timestamp
	var lastExecution models.Execution
	if err := s.db.Where("job_id = ?", job.ID).
		Order("timestamp DESC").
		First(&lastExecution).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	} else {
		job.LastExecution = &lastExecution.Timestamp
	}

	return nil
}
