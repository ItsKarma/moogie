package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/itskarma/moogie/api/internal/handlers"
	"github.com/itskarma/moogie/api/internal/services"
	"github.com/itskarma/moogie/api/internal/websocket"
	"github.com/itskarma/moogie/api/pkg/config"
	"github.com/itskarma/moogie/api/pkg/database"
)

// @title Moogie API
// @version 1.0
// @description API for Moogie monitoring dashboard
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db := database.Connect(cfg)

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Initialize services
	jobService := services.NewJobService(db)
	executionService := services.NewExecutionService(db, jobService)
	dashboardService := services.NewDashboardService(db, jobService, executionService)

	// Initialize handlers
	handler := handlers.NewHandler(jobService, executionService, dashboardService, wsHub)

	// Setup Gin router
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Setup CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.AllowedOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(corsConfig))

	// Setup routes
	setupRoutes(router, handler)

	// Start server
	log.Printf("Starting server on port %s", cfg.AppPort)
	log.Printf("Environment: %s", cfg.AppEnv)
	log.Printf("Allowed origins: %v", cfg.AllowedOrigins)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(router *gin.Engine, handler *handlers.Handler) {
	// Health check
	router.GET("/health", handler.HealthCheck)

	// WebSocket endpoint
	router.GET("/ws", handler.HandleWebSocket)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Jobs
		jobs := v1.Group("/jobs")
		{
			jobs.GET("", handler.GetJobs)
			jobs.GET("/:id", handler.GetJob)
		}

		// Executions
		executions := v1.Group("/executions")
		{
			executions.POST("", handler.CreateExecution)
		}

		// Dashboard
		dashboard := v1.Group("/dashboard")
		{
			dashboard.GET("/summary", handler.GetDashboardSummary)
		}
	}
}
