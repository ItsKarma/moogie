package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/itskarma/moogie/api/internal/models"
	"gorm.io/gorm"
)

type DashboardService struct {
	db               *gorm.DB
	jobService       *JobService
	executionService *ExecutionService
}

func NewDashboardService(db *gorm.DB, jobService *JobService, executionService *ExecutionService) *DashboardService {
	return &DashboardService{
		db:               db,
		jobService:       jobService,
		executionService: executionService,
	}
}

// GetSummary returns aggregated dashboard metrics
func (s *DashboardService) GetSummary(from, to time.Time) (*models.DashboardSummary, error) {
	summary := &models.DashboardSummary{}

	// Get total jobs count
	if err := s.db.Model(&models.Job{}).Count(&summary.TotalJobs).Error; err != nil {
		return nil, fmt.Errorf("failed to count total jobs: %w", err)
	}

	// Get active jobs count
	if err := s.db.Model(&models.Job{}).Where("enabled = ?", true).Count(&summary.ActiveJobs).Error; err != nil {
		return nil, fmt.Errorf("failed to count active jobs: %w", err)
	}

	// Get total executions in date range
	if err := s.db.Model(&models.Execution{}).
		Where("timestamp BETWEEN ? AND ?", from, to).
		Count(&summary.TotalExecutions).Error; err != nil {
		return nil, fmt.Errorf("failed to count total executions: %w", err)
	}

	// Calculate overall success rate
	if summary.TotalExecutions > 0 {
		var successfulExecutions int64
		if err := s.db.Model(&models.Execution{}).
			Where("status = 'success' AND timestamp BETWEEN ? AND ?", from, to).
			Count(&successfulExecutions).Error; err != nil {
			return nil, fmt.Errorf("failed to count successful executions: %w", err)
		}
		summary.OverallSuccess = float64(successfulExecutions) / float64(summary.TotalExecutions) * 100
	}

	// Get job summaries
	jobSummaries, err := s.getJobSummaries(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get job summaries: %w", err)
	}
	summary.JobSummaries = jobSummaries

	// Get recent activity
	recentActivity, err := s.executionService.GetRecentExecutions(20)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}
	summary.RecentActivity = recentActivity

	// Get status breakdown
	statusBreakdown, err := s.getStatusBreakdown(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get status breakdown: %w", err)
	}
	summary.StatusBreakdown = statusBreakdown

	// Get type breakdown
	typeBreakdown, err := s.getTypeBreakdown()
	if err != nil {
		return nil, fmt.Errorf("failed to get type breakdown: %w", err)
	}
	summary.TypeBreakdown = typeBreakdown

	return summary, nil
}

// getJobSummaries returns summary metrics for each job
func (s *DashboardService) getJobSummaries(from, to time.Time) ([]models.JobSummary, error) {
	var jobs []models.Job
	
	// Load all jobs without preloading executions
	if err := s.db.Find(&jobs).Error; err != nil {
		return nil, err
	}

	var summaries []models.JobSummary
	for _, job := range jobs {
		// Compute metrics for this job within date range
		if err := s.jobService.computeJobMetrics(&job, from, to); err != nil {
			return nil, err
		}

		// Count executions for this job in the date range
		var executionCount int64
		if err := s.db.Model(&models.Execution{}).
			Where("job_id = ? AND timestamp BETWEEN ? AND ?", job.ID, from, to).
			Count(&executionCount).Error; err != nil {
			return nil, err
		}

		// Get last 10 executions for this job (regardless of date range)
		var recentExecutions []models.Execution
		if err := s.db.Where("job_id = ?", job.ID).
			Select("id", "job_id", "status", "response_time", "details", "timestamp").
			Order("timestamp DESC").
			Limit(10).
			Find(&recentExecutions).Error; err != nil {
			return nil, err
		}

		// Extract labels from config
		labels := extractLabels(job.Config)

		summary := models.JobSummary{
			ID:               job.ID,
			Name:             job.Name,
			Type:             job.Type,
			Enabled:          job.Enabled,
			SuccessRate:      job.SuccessRate,
			LastExecution:    job.LastExecution,
			AvgResponseTime:  job.AvgResponseTime,
			ExecutionCount:   executionCount,
			RecentExecutions: recentExecutions, // Last 10 executions (always, regardless of date range)
			Labels:           labels,
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// extractLabels extracts labels from job config JSON
func extractLabels(configJSON []byte) models.Labels {
	var config struct {
		Metadata struct {
			Labels models.Labels `json:"labels"`
		} `json:"metadata"`
	}

	// Try to unmarshal the config
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return models.Labels{} // Return empty labels if parsing fails
	}

	return config.Metadata.Labels
}

// getStatusBreakdown returns count of executions by status
func (s *DashboardService) getStatusBreakdown(from, to time.Time) (map[string]int64, error) {
	type statusCount struct {
		Status string
		Count  int64
	}

	var results []statusCount
	if err := s.db.Model(&models.Execution{}).
		Select("status, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", from, to).
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	breakdown := make(map[string]int64)
	for _, result := range results {
		breakdown[result.Status] = result.Count
	}

	return breakdown, nil
}

// getTypeBreakdown returns count of jobs by type
func (s *DashboardService) getTypeBreakdown() (map[string]int64, error) {
	type typeCount struct {
		Type  string
		Count int64
	}

	var results []typeCount
	if err := s.db.Model(&models.Job{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	breakdown := make(map[string]int64)
	for _, result := range results {
		breakdown[result.Type] = result.Count
	}

	return breakdown, nil
}
