// internal/api/router.go
package api

import (
	"backend/internal/api/handlers"
	"backend/internal/api/middleware"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all API routes
func SetupRouter(
	projectService *services.ProjectService,
	targetService *services.TargetService,
	scanService *services.ScanService,
	findingService *services.FindingService,
	queueService *services.QueueService,
) *gin.Engine {
	// Create router with default logger and recovery middleware
	router := gin.Default()

	// Add custom middleware
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.ErrorHandler())

	// Optional rate limiter - uncomment if needed
	// router.Use(middleware.RateLimiter())

	// Create handlers
	projectHandler := handlers.NewProjectHandler(projectService, targetService)
	targetHandler := handlers.NewTargetHandler(targetService)
	scanHandler := handlers.NewScanHandler(scanService, queueService, projectService, targetService)
	findingHandler := handlers.NewFindingHandler(findingService)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Projects
		projects := v1.Group("/projects")
		{
			projects.GET("", projectHandler.GetProjects)
			projects.POST("", projectHandler.CreateProject)
			projects.GET("/:id", projectHandler.GetProject)
			projects.PUT("/:id", projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)
			projects.GET("/:id/targets", projectHandler.GetProjectTargets)
			projects.GET("/:id/scans", projectHandler.GetProjectScans)
			projects.GET("/:id/findings", projectHandler.GetProjectFindings)
		}

		// Targets
		targets := v1.Group("/targets")
		{
			targets.GET("", targetHandler.GetTargets)
			targets.POST("", targetHandler.CreateTarget)
			targets.POST("/bulk", targetHandler.BulkCreateTargets)
			targets.GET("/:id", targetHandler.GetTarget)
			targets.PUT("/:id", targetHandler.UpdateTarget)
			targets.DELETE("/:id", targetHandler.DeleteTarget)
			targets.GET("/:id/findings", targetHandler.GetTargetFindings)
		}

		// Scans
		scans := v1.Group("/scans")
		{
			scans.GET("", scanHandler.GetScans)
			scans.POST("", scanHandler.StartScan)
			scans.GET("/:id", scanHandler.GetScan)
			scans.POST("/:id/cancel", scanHandler.CancelScan)
			scans.GET("/:id/findings", scanHandler.GetScanFindings)
			scans.GET("/:id/tasks", scanHandler.GetScanTasks)
		}

		// Scan configurations
		scanConfigs := v1.Group("/scan-configs")
		{
			scanConfigs.GET("", scanHandler.GetScanConfigs)
			scanConfigs.POST("", scanHandler.CreateScanConfig)
			scanConfigs.PUT("/:id", scanHandler.UpdateScanConfig)
			scanConfigs.DELETE("/:id", scanHandler.DeleteScanConfig)
		}

		// Findings
		findings := v1.Group("/findings")
		{
			findings.GET("", findingHandler.GetFindings)
			findings.POST("", findingHandler.CreateFinding)
			findings.GET("/latest", findingHandler.GetLatestFindings)
			findings.GET("/summary/:project_id", findingHandler.GetFindingsSummary)
			findings.GET("/:id", findingHandler.GetFinding)
			findings.PUT("/:id", findingHandler.UpdateFinding)
			findings.DELETE("/:id", findingHandler.DeleteFinding)
			findings.POST("/bulk-update", findingHandler.BulkUpdateFindings)
			findings.PUT("/:id/fixed/:fixed", findingHandler.MarkFixed)
			findings.PUT("/:id/verified/:verified", findingHandler.MarkVerified)
		}
	}

	// Serve API documentation if in development mode
	// Uncomment if you're using Swagger or similar
	/*
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	*/

	return router
}
