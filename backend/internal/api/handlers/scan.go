package handlers

import (
	"net/http"
	"time"

	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ScanHandler struct {
	scanService    *services.ScanService
	queueService   *services.QueueService
	projectService *services.ProjectService
	targetService  *services.TargetService
}

func NewScanHandler(
	scanService *services.ScanService,
	queueService *services.QueueService,
	projectService *services.ProjectService,
	targetService *services.TargetService,
) *ScanHandler {
	return &ScanHandler{
		scanService:    scanService,
		queueService:   queueService,
		projectService: projectService,
		targetService:  targetService,
	}
}

// GetScans returns all scans
// @Summary Get all scans
// @Description Get all scans across all projects
// @Tags scans
// @Accept json
// @Produce json
// @Success 200 {array} models.Scan
// @Failure 500 {object} map[string]string
// @Router /api/v1/scans [get]
func (h *ScanHandler) GetScans(c *gin.Context) {
	scans, err := h.scanService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve scans"})
		return
	}

	c.JSON(http.StatusOK, scans)
}

// GetScan returns a specific scan by ID
// @Summary Get a scan
// @Description Get a specific scan by ID
// @Tags scans
// @Accept json
// @Produce json
// @Param id path string true "Scan ID"
// @Success 200 {object} models.Scan
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/scans/{id} [get]
func (h *ScanHandler) GetScan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scan ID format"})
		return
	}

	scan, err := h.scanService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scan not found"})
		return
	}

	c.JSON(http.StatusOK, scan)
}

// StartScan creates a new scan and queues it
// @Summary Start a new scan
// @Description Create a new scan and queue it
// @Tags scans
// @Accept json
// @Produce json
// @Param scan body models.StartScanInput true "Scan Details"
// @Success 202 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/scans [post]
func (h *ScanHandler) StartScan(c *gin.Context) {
	var input models.StartScanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify project exists
	project, err := h.projectService.GetByID(input.ProjectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Verify scan config exists
	scanConfig, err := h.scanService.GetScanConfigByID(input.ScanConfigID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scan configuration not found"})
		return
	}

	// Create scan record
	scan := &models.Scan{
		ProjectID:    project.ID,
		ScanConfigID: scanConfig.ID,
		Status:       models.StatusPending,
		CreatedAt:    time.Now(),
	}

	err = h.scanService.Create(scan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create scan"})
		return
	}

	// Determine targets
	var targets []models.Target
	if len(input.TargetIDs) > 0 {
		// Use specified targets
		for _, targetID := range input.TargetIDs {
			target, err := h.targetService.GetByID(targetID)
			if err == nil {
				targets = append(targets, *target)
			}
		}
	} else {
		// Use all targets from the project
		projectTargets, err := h.targetService.GetByProjectID(project.ID)
		if err == nil {
			targets = projectTargets
		}
	}

	if len(targets) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid targets found for scanning"})
		return
	}

	// Queue scan tasks based on scanner type
	scanRequest := services.ScanRequest{
		ScanID:      scan.ID,
		ScannerType: scanConfig.ScannerType,
		Targets:     targets,
		Parameters:  scanConfig.Parameters,
	}

	err = h.queueService.QueueScan(scanRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue scan"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Scan queued successfully",
		"scan_id": scan.ID,
	})
}

// CancelScan cancels a running scan
// @Summary Cancel a scan
// @Description Cancel a running scan
// @Tags scans
// @Accept json
// @Produce json
// @Param id path string true "Scan ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/scans/{id}/cancel [post]
func (h *ScanHandler) CancelScan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scan ID format"})
		return
	}

	scan, err := h.scanService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scan not found"})
		return
	}

	if scan.Status != models.StatusPending && scan.Status != models.StatusRunning {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can only cancel pending or running scans"})
		return
	}

	// Update scan status
	err = h.scanService.UpdateStatus(scan.ID, models.StatusCancelled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update scan status"})
		return
	}

	// Send cancel message to queue
	err = h.queueService.CancelScan(scan.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel scan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scan cancelled successfully"})
}

// GetScanFindings returns all findings for a specific scan
// @Summary Get scan findings
// @Description Get all findings for a specific scan
// @Tags scans
// @Accept json
// @Produce json
// @Param id path string true "Scan ID"
// @Success 200 {array} models.Finding
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/scans/{id}/findings [get]
func (h *ScanHandler) GetScanFindings(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scan ID format"})
		return
	}

	findings, err := h.scanService.GetFindings(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve findings"})
		return
	}

	c.JSON(http.StatusOK, findings)
}

// GetScanTasks returns all tasks for a specific scan
// @Summary Get scan tasks
// @Description Get all tasks for a specific scan
// @Tags scans
// @Accept json
// @Produce json
// @Param id path string true "Scan ID"
// @Success 200 {array} models.ScanTask
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/scans/{id}/tasks [get]
func (h *ScanHandler) GetScanTasks(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scan ID format"})
		return
	}

	tasks, err := h.scanService.GetScanTasks(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve scan tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetScanConfigs returns all scan configurations
// @Summary Get scan configurations
// @Description Get all scan configurations
// @Tags scan-configs
// @Accept json
// @Produce json
// @Success 200 {array} models.ScanConfig
// @Failure 500 {object} map[string]string
// @Router /api/v1/scan-configs [get]
func (h *ScanHandler) GetScanConfigs(c *gin.Context) {
	configs, err := h.scanService.GetAllScanConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve scan configurations"})
		return
	}

	c.JSON(http.StatusOK, configs)
}

// GetScanConfig returns all scan configurations
// @Summary Get scan configurations
// @Description Get all scan configurations
// @Tags scan-configs
// @Accept json
// @Produce json
// @Param id path string true "Scan Config ID"
// @Success 200 {array} models.ScanConfig
// @Failure 500 {object} map[string]string
// @Router /api/v1/scan-configs [get]
func (h *ScanHandler) GetScanConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scan ID format"})
		return
	}

	config, err := h.scanService.GetScanConfigByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve scan configuration"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// CreateScanConfig creates a new scan configuration
// @Summary Create a scan configuration
// @Description Create a new scan configuration
// @Tags scan-configs
// @Accept json
// @Produce json
// @Param config body object true "Scan Configuration Details"
// @Success 201 {object} models.ScanConfig
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/scan-configs [post]
func (h *ScanHandler) CreateScanConfig(c *gin.Context) {
	var input struct {
		Name        string       `json:"name" binding:"required"`
		ScannerType string       `json:"scanner_type" binding:"required"`
		Parameters  models.JSONB `json:"parameters"`
		Active      bool         `json:"active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := &models.ScanConfig{
		Name:        input.Name,
		ScannerType: input.ScannerType,
		Parameters:  input.Parameters,
		Active:      input.Active,
	}

	err := h.scanService.CreateScanConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create scan configuration"})
		return
	}

	c.JSON(http.StatusCreated, config)
}

// UpdateScanConfig updates an existing scan configuration
// @Summary Update a scan configuration
// @Description Update an existing scan configuration
// @Tags scan-configs
// @Accept json
// @Produce json
// @Param id path string true "Scan Config ID"
// @Param config body object true "Updated Scan Configuration Details"
// @Success 200 {object} models.ScanConfig
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/scan-configs/{id} [put]
func (h *ScanHandler) UpdateScanConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scan config ID format"})
		return
	}

	var input struct {
		Name        string       `json:"name"`
		ScannerType string       `json:"scanner_type"`
		Parameters  models.JSONB `json:"parameters"`
		Active      *bool        `json:"active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config, err := h.scanService.GetScanConfigByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scan configuration not found"})
		return
	}

	// Update fields if provided
	if input.Name != "" {
		config.Name = input.Name
	}

	if input.ScannerType != "" {
		config.ScannerType = input.ScannerType
	}

	if input.Parameters != nil {
		config.Parameters = input.Parameters
	}

	if input.Active != nil {
		config.Active = *input.Active
	}

	err = h.scanService.UpdateScanConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update scan configuration"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// DeleteScanConfig deletes a scan configuration
// @Summary Delete a scan configuration
// @Description Delete a scan configuration
// @Tags scan-configs
// @Accept json
// @Produce json
// @Param id path string true "Scan Config ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/scan-configs/{id} [delete]
func (h *ScanHandler) DeleteScanConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scan config ID format"})
		return
	}

	err = h.scanService.DeleteScanConfig(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete scan configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scan configuration deleted successfully"})
}
