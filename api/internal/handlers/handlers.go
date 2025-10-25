package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itskarma/moogie/api/internal/models"
	"github.com/itskarma/moogie/api/internal/services"
	"github.com/itskarma/moogie/api/internal/websocket"
)

type Handler struct {
	jobService       *services.JobService
	executionService *services.ExecutionService
	dashboardService *services.DashboardService
	wsHub            *websocket.Hub
}

// NewHandler creates a new handler instance
func NewHandler(
	jobService *services.JobService,
	executionService *services.ExecutionService,
	dashboardService *services.DashboardService,
	wsHub *websocket.Hub,
) *Handler {
	return &Handler{
		jobService:       jobService,
		executionService: executionService,
		dashboardService: dashboardService,
		wsHub:            wsHub,
	}
}

// @Summary Get all jobs
// @Description Get all monitoring jobs with their success rates and recent metrics
// @Tags jobs
// @Accept json
// @Produce json
// @Param from query string false "Start date (YYYY-MM-DD format)"
// @Param to query string false "End date (YYYY-MM-DD format)"
// @Success 200 {array} models.Job
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs [get]
func (h *Handler) GetJobs(c *gin.Context) {
	// Parse date range from query parameters
	from, to, err := parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobs, err := h.jobService.GetAllJobs(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

// @Summary Get job by ID
// @Description Get a specific job with its execution history
// @Tags jobs
// @Accept json
// @Produce json
// @Param id path int true "Job ID"
// @Param from query string false "Start date for executions (YYYY-MM-DD format)"
// @Param to query string false "End date for executions (YYYY-MM-DD format)"
// @Param limit query int false "Limit number of executions returned" default(100)
// @Success 200 {object} models.Job
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/{id} [get]
func (h *Handler) GetJob(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	// Parse date range and limit
	from, to, err := parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := 100 // default
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	job, err := h.jobService.GetJobByID(uint(id), from, to, limit)
	if err != nil {
		if err.Error() == "job not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, job)
}

// @Summary Create execution result
// @Description Create a new execution result (called by runner service)
// @Tags executions
// @Accept json
// @Produce json
// @Param execution body models.CreateExecutionRequest true "Execution data"
// @Success 201 {object} models.Execution
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /executions [post]
func (h *Handler) CreateExecution(c *gin.Context) {
	var req models.CreateExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	execution, err := h.executionService.CreateExecution(&req)
	if err != nil {
		if err.Error() == "job not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Broadcast the new execution to WebSocket clients
	h.wsHub.BroadcastExecutionCreated(execution)

	c.JSON(http.StatusCreated, execution)
}

// @Summary Get dashboard summary
// @Description Get aggregated dashboard metrics and recent activity
// @Tags dashboard
// @Accept json
// @Produce json
// @Param from query string false "Start date (YYYY-MM-DD format)"
// @Param to query string false "End date (YYYY-MM-DD format)"
// @Success 200 {object} models.DashboardSummary
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dashboard/summary [get]
func (h *Handler) GetDashboardSummary(c *gin.Context) {
	// Parse date range from query parameters
	from, to, err := parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	summary, err := h.dashboardService.GetSummary(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// @Summary WebSocket endpoint
// @Description WebSocket endpoint for real-time updates
// @Tags websocket
// @Router /ws [get]
func (h *Handler) HandleWebSocket(c *gin.Context) {
	h.wsHub.HandleWebSocketConnection(c)
}

// @Summary Health check
// @Description Health check endpoint
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// parseDateRange parses the from and to query parameters
// Returns default range (last 7 days) if not provided
func parseDateRange(c *gin.Context) (time.Time, time.Time, error) {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to time.Time
	var err error

	if fromStr == "" || toStr == "" {
		// Default to last 7 days if not provided
		to = time.Now()
		from = to.AddDate(0, 0, -7)
	} else {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid from date format, expected YYYY-MM-DD")
		}

		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid to date format, expected YYYY-MM-DD")
		}

		// Set to end of day for 'to' date
		to = to.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	return from, to, nil
}
